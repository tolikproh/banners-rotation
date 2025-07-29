package internalhttp

import "net/http"

func (s *Server) initRoute() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(helloHandler))

	mux.Handle("GET /banner", s.serveHandler(s.getBanner))
	mux.Handle("POST /banner", s.serveHandler(s.addBannerToSlot))
	mux.Handle("POST /click", s.serveHandler(s.clickBanner))
	mux.Handle("DELETE /banner", s.serveHandler(s.deleteBannerFromSlot))

	mux.Handle("POST /create/banner", s.serveHandler(s.createBanner))
	mux.Handle("POST /create/slot", s.serveHandler(s.createSlot))
	mux.Handle("POST /create/group", s.serveHandler(s.createSocislGroup))

	mux.Handle("DELETE /delete/banner", s.serveHandler(s.deleteBanner))
	mux.Handle("DELETE /delete/slot", s.serveHandler(s.deleteSlot))
	mux.Handle("DELETE /delete/group", s.serveHandler(s.deleteSocialGroup))

	var handler http.Handler = mux
	handler = s.loggingMiddleware(handler)

	return handler
}
