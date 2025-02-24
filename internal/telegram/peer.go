package telegram

import "github.com/gotd/td/tg"

// GetInputPeer фильтрует чаты и пользователей (игнорируя каналы)
func GetInputPeer(peer tg.PeerClass) tg.InputPeerClass {
	switch p := peer.(type) {
	case *tg.PeerUser:
		return &tg.InputPeerUser{UserID: p.UserID}
	case *tg.PeerChat:
		return &tg.InputPeerChat{ChatID: p.ChatID}
	default:
		return nil // Исключаем каналы
	}
}

