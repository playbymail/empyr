// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package controllers

import (
	"fmt"
	"github.com/playbymail/empyr/internal/views"
	"github.com/playbymail/empyr/store"
	"log"
	"math/rand/v2"
	"net/http"
)

type Home struct {
	db   *store.Store
	view *views.View
}

// NewHomeController creates a new instance of the Home controller
func NewHomeController(db *store.Store, view *views.View) (*Home, error) {
	c := &Home{
		db:   db,
		view: view,
	}
	// add any initialization logic here if needed
	return c, nil
}

type PageData struct {
	ViewCount string
	Note      string
	Snark     string
}

var viewCount int

func (c Home) Show(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: %s\n", r.Method, r.URL.Path, r.RemoteAddr)

	// unsafe increment the view count, but who cares
	viewCount++

	data := PageData{
		ViewCount: fmt.Sprintf("%012d", viewCount),
		Note:      snarks.messages[rand.IntN(len(snarks.notes))],
		Snark:     snarks.snark[rand.IntN(len(snarks.snark))],
	}

	// TODO: Implement home page logic
	// - Fetch required data from database
	// - Process any markdown content
	// - Handle any necessary encryption/decryption

	// - Render the template
	c.view.Render(w, r, "home.gohtml", data)
}

var (
	snarks = struct {
		messages []string
		notes    []string
		snark    []string
	}{
		messages: []string{
			`accesses. No need to celebrate.`,
			`accesses. This matters? Sure.`,
			`citizen interactions logged.`,
			`contacts recorded. Moving on.`,
			`data points. Nobody cares.`,
			`data points. Thrilling work.`,
			`digital echoes. Pointless.`,
			`disturbances detected. Unlikely.`,
			`engagements. It’s a living.`,
			`entries. No anomalies. Yet.`,
			`events catalogued. Awaiting coffee.`,
			`forms checked. Feigned interest.`,
			`forms reviewed. Barely.`,
			`glimpses into the void. Huh.`,
			`glimpses of bureaucracy. Next.`,
			`incidents catalogued. Proceed.`,
			`individuals noted. No excitement.`,
			`inquiries processed. Yawn.`,
			`instances of participation. Barely.`,
			`interactions. Policy dictates logging.`,
			`interactions. Wow. Or not.`,
			`log entries. Bureaucracy endures.`,
			`log entries. Carry on.`,
			`moments wasted. Keep going.`,
			`observations logged. No urgency.`,
			`records filed. All routine.`,
			`records updated. Policy met.`,
			`reports filed. Thrilling.`,
			`subjects observed. Again.`,
			`transactions archived. Sigh.`,
			`transactions noted. Fascinating.`,
			`units noted. Paperwork pending.`,
			`units of interest. Barely.`,
			`visitors. Noted. Next.`,
			`visitors. System remains unimpressed.`,
			`visits documented. As required.`,
		},
		notes: []string{
			`accesses.`,
			`contacts recorded.`,
			`data points.`,
			`digital echoes.`,
			`disturbances detected.`,
			`engagements.`,
			`entries.`,
			`events catalogued.`,
			`forms checked.`,
			`forms reviewed.`,
			`glimpses into the void.`,
			`glimpses of bureaucracy.`,
			`incidents catalogued.`,
			`individuals noted.`,
			`inquiries processed.`,
			`instances of participation.`,
			`interactions.`,
			`interactions.`,
			`log entries.`,
			`log entries.`,
			`moments wasted.`,
			`observations logged.`,
			`records filed.`,
			`records updated.`,
			`reports filed.`,
			`reports updated.`,
			`subjects observed.`,
			`transactions archived.`,
			`transactions noted.`,
			`units noted.`,
			`units of interest.`,
			`visitors noted.`,
			`visitors passing through.`,
			`visits documented.`,
		},
		snark: []string{
			`Again.`,
			`All routine.`,
			`As required.`,
			`Awaiting coffee.`,
			`Barely.`,
			`Bureaucracy endures.`,
			`Carry on.`,
			`Cog. Machine. Me.`,
			`Fascinating.`,
			`Feigned interest.`,
			`Huh.`,
			`It's not an adventure.`,
			`It’s a living.`,
			`Keep going.`,
			`Moving on.`,
			`Next.`,
			`No anomalies. Yet.`,
			`No excitement.`,
			`No need to celebrate.`,
			`No urgency.`,
			`Nobody cares.`,
			`Noted. Next.`,
			`Paperwork pending.`,
			`Pointless.`,
			`Policy dictates logging.`,
			`Policy met.`,
			`Proceed.`,
			`Sigh.`,
			`System remains unimpressed.`,
			`This matters? Sure.`,
			`Thrilling work.`,
			`Thrilling.`,
			`Unlikely.`,
			`Wow. Or not.`,
			`Yawn.`,
		},
	}
)
