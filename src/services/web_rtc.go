package services

import (
	"fmt"
	"sync"

	"github.com/pion/webrtc/v4"
)
type Participant struct {
	ID int
	PC *webrtc.PeerConnection
}

type RTCSession struct {
	PeerConnection *webrtc.PeerConnection
	Participants map[int]*Participant
	Mutex sync.Mutex
}

type RTCService struct {
	Sessions map[int]*RTCSession
	Mutex sync.Mutex
	OnICECandidate func(meetingID, userID int, candidate string)
}

func NewRTCService() *RTCService {
	return &RTCService{
		Sessions: make(map[int]*RTCSession),
	}
}

func (service *RTCService) CreateSession(meetingID, userID int) (string, error) {
	service.Mutex.Lock()
	defer service.Mutex.Unlock()

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return "", err
	}

	session := &RTCSession{
		PeerConnection: pc,
		Participants: make(map[int]*Participant),
	}
	service.Sessions[meetingID] = session

	_, err = pc.CreateDataChannel("chat", nil)
	if err != nil {
		return "", err
	}

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		return "", err
	}

	err = pc.SetLocalDescription(offer)
	if err != nil {
		return "", err
	}

	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil && service.OnICECandidate != nil {
			service.OnICECandidate(meetingID, userID, c.ToJSON().Candidate)
		}
	})

	return offer.SDP, nil
}

func (service *RTCService) JoinSession(meetingID, userID int, offerSDP string) (string, error) {
	service.Mutex.Lock()
	session, ok := service.Sessions[meetingID]
	service.Mutex.Unlock()
	if !ok {
		return "", fmt.Errorf("session not found")
	}

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return "", err
	}

	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	newParticipant := &Participant{ID: userID, PC: pc}
	session.Participants[userID] = newParticipant

	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		session.Mutex.Lock()

		for _, other := range session.Participants {
			if other.ID == userID {
				continue
			}

			localTrack, err := webrtc.NewTrackLocalStaticRTP(track.Codec().RTPCodecCapability, track.ID(), track.StreamID())
			if err != nil {
				fmt.Println("error creating localStrack:", err)
				continue
			}

			other.PC.AddTrack(localTrack)

			go func() {
				buf := make([]byte, 1500)
				for {
					n, _, readErr := track.Read(buf)
					if readErr != nil {
						return 
					}
					_, writeErr := localTrack.Write(buf[:n])
					if writeErr != nil {
						return 
					}
				}
			}()
		}
		session.Mutex.Unlock()
	})

	err = pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP: offerSDP,
	})
	if err != nil {
		return "", err
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return "", err
	}

	err = pc.SetLocalDescription(answer)
	if err != nil {
		return "", err
	}

	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil && service.OnICECandidate != nil {
			service.OnICECandidate(meetingID, userID, c.ToJSON().Candidate)
		}
	})

	return answer.SDP, nil
}

func (service *RTCService) LeaveSession(meetingID, userID int) {
	service.Mutex.Lock()
	session, ok := service.Sessions[meetingID]
	service.Mutex.Unlock()
	if !ok {
		return 
	}

	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	if p, exists := session.Participants[userID]; exists {
		p.PC.Close()
		delete(session.Participants, userID)
	}
}

func (service *RTCService) AddIceCandidate(meetingID int, candidate string, userID int) error {
	service.Mutex.Lock()
	session, ok := service.Sessions[meetingID]
	service.Mutex.Unlock()
	if !ok {
		return fmt.Errorf("session not found")
	}

	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	p, exists := session.Participants[userID]
	if !exists {
		return fmt.Errorf("participant not found")
	}

	return p.PC.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate})
}