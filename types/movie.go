/*
 * Copyright (c) 2024.  Jan Eike Suchard. All rights reserved
 * SPDX-License-Identifier: MIT
 */

package types

import "github.com/jackc/pgx/v5/pgtype"

// Movie represents a single movie screened in the cinema. A movie can have
// multiple ScreeningTimes
type Movie struct {
	// ID contains the unique id of the movie
	ID pgtype.UUID `json:"id" db:"id"`

	// Title contains the german title of the movie
	Title pgtype.Text `json:"title" db:"title"`

	// OriginalTitle contains the title in the original movies language
	OriginalTitle pgtype.Text `json:"originalTitle" db:"originalTitle"`

	// Description contains a long textual description of the movie. It may
	// contain markdown which needs to be rendered on the frontend side
	Description pgtype.Text `json:"description" db:"description"`

	// ScreeningTimes is an array containing the ScreeningTime instances for
	// this movie
	ScreeningTimes []ScreeningTime `json:"screeningTimes" db:"screening_times"`

	// AudioLanguage contains the ISO 639 (Set 1) code of the spoken language
	// for this movie
	AudioLanguage pgtype.Text `json:"audioLanguage" db:"audio_language"`

	// SubtitleLanguage contains the ISO 639 (Set 1) code of the spoken
	// language for this movie
	SubtitleLanguage pgtype.Text `json:"subtitleLanguage" db:"subtitle_language"`

	// Duration contains the running length of the movie in minutes
	Duration int `json:"duration" db:"duration"`

	// AdditionalInformation contains some more, but not specified information
	// about the movie
	AdditionalInformation interface{} `json:"additionalInformation" db:"additional_information"`
}
