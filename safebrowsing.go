// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"os"
	"time"

	"github.com/flosch/pongo2"
	"github.com/google/safebrowsing"
)

var safeBrowser *safebrowsing.SafeBrowser

func init() {
	pongo2.RegisterFilter("threatdefinition", safeTypeToStringFilter)
}

func initSafeBrowsing() {
	if conf.SafeBrowsing.APIKey == "" {
		return
	}

	debug.Println("safebrowsing support enabled, initializing")

	// Validate the part of the config that we can.
	if conf.SafeBrowsing.UpdatePeriod < 30*time.Minute {
		// Minimum 30m.
		conf.SafeBrowsing.UpdatePeriod = 30 * time.Minute
	}
	if conf.SafeBrowsing.UpdatePeriod > 168*time.Hour {
		// Maximum 7 days.
		conf.SafeBrowsing.UpdatePeriod = 168 * time.Minute
	}

	var err error
	safeBrowser, err = safebrowsing.NewSafeBrowser(safebrowsing.Config{
		APIKey:         conf.SafeBrowsing.APIKey,
		DBPath:         conf.SafeBrowsing.DBPath,
		UpdatePeriod:   conf.SafeBrowsing.UpdatePeriod,
		RequestTimeout: 15 * time.Second,
		Logger:         os.Stdout,
	})
	if err != nil {
		debug.Fatalf("error initializing google safebrowsing: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if err = safeBrowser.WaitUntilReady(ctx); err != nil {
		debug.Fatalf("error initializing google safebrowsing: %v", err)
	}
}

func safeTypeToString(t safebrowsing.ThreatType) string {
	switch t {
	case safebrowsing.ThreatType_Malware:
		return "Site is known for hosting malware"
	case safebrowsing.ThreatType_PotentiallyHarmfulApplication:
		return "Site provides potentially harmful applications"
	case safebrowsing.ThreatType_SocialEngineering:
		return "Site is known for social engineering"
	case safebrowsing.ThreatType_UnwantedSoftware:
		return "Site provides unwanted software"
	}

	return "Unknown threat"
}

func safeTypeToStringFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	input := in.Integer()
	t := safebrowsing.ThreatType(input)
	return pongo2.AsValue(safeTypeToString(t)), nil
}
