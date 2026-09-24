# Podcast API Go Library

[![Go CI](https://github.com/ListenNotes/podcast-api-go/actions/workflows/go.yml/badge.svg)](https://github.com/ListenNotes/podcast-api-go/actions/workflows/go.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/ListenNotes/podcast-api-go/v3.svg)](https://pkg.go.dev/github.com/ListenNotes/podcast-api-go/v3)

The official Go client for the [Listen Notes Podcast API](https://www.listennotes.com/api/).
Search podcasts and episodes, fetch metadata, and manage playlists.
Questions: [hello@listennotes.com](mailto:hello@listennotes.com).

## Installation

Requires Go 1.26 or newer. The SDK uses only the Go standard library.

```sh
go get github.com/ListenNotes/podcast-api-go/v3
```

## Usage

Path identifiers are positional strings; other fields use `map[string]string`.
An empty API key selects the stateless public mock. To call production, pass your
[Listen API key](https://www.listennotes.com/api/dashboard/#apps) to `NewClient`.

```go
package main

import (
	"fmt"

	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.Search(map[string]string{"q": "star wars", "type": "episode"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
	fmt.Println(response.Stats.Usage, response.Stats.FreeQuota)
}
```

Reuse clients for connection pooling. Each client retains its own API key and
configuration. Defaults: 30-second total timeout, 10-second connection timeout,
and redirects disabled. The SDK does not add retries. `WithHTTPClient` accepts
your own HTTP client and its timeout, redirect, proxy, and transport policies.
`WithBaseURL` supports local test servers and sends your key to that chosen URL.

### Upgrading from 1.x

- Change imports and `go get` commands to
  `github.com/ListenNotes/podcast-api-go/v3`, as required by Go module versioning.
- Use Go 1.26+. Existing method names, positional identifiers, string parameter
  maps, `Response.Data`, `Response.Stats`, and `ToJSON()` remain available.
- The `HTTPClient` interface adds five playlist methods; custom implementations
  of that interface must add them too.
- Match errors with `errors.Is(err, listennotes.ErrUnauthorized)`, replacing
  direct equality checks. Use `errors.As(err, &apiError)` with
  `var apiError *listennotes.APIError` for status, headers, and response body.
  HTTP 403 is `ErrForbidden`; all non-2xx statuses return errors.
- Responses also expose `StatusCode`, `Headers`, and the unmodified `Body`,
  including on HTTP errors. Use named fields when constructing response literals.
- POST/PUT fields are sent in the form body; only declared query fields enter the
  query string. Path segments are escaped. Empty notes and descriptions are
  preserved so they can be cleared. DELETE parameters such as `reason` remain
  query parameters. Maps passed by callers are not modified.
- Default clients no longer follow redirects. Statistics parsing remains best
  effort and cannot turn a successful write into an error.

## Development

Go's built-in module, formatting, testing, and vet tools cover this project.
There are no third-party dependencies or separate tool installations.

```sh
go test -race -count=1 ./...
go vet ./...
test -z "$(gofmt -l .)"
go mod tidy -diff
go build ./...
```

Default tests use loopback HTTP fixtures or mocked transports; they never send
requests to production or the public mock. Generated README examples compile
as Go examples without executing requests. Run the public mock suite separately:

```sh
go test -tags=integration -run '^TestMockIntegration' -count=1 -timeout=3m ./...
```

Integration tests use only the fixed public mock origin, with no API key,
environment-selected destination, proxy, or redirects. They cover all 30 methods.
The mock is stateless; these tests do not establish production write persistence
or authorization. `go run ./example` makes one explicit read-only mock request.

The monorepo's `sync.py go` generates methods, the contract, test dispatch,
compile-only examples, and the marked README sections. Do not hand-edit generated
outputs. This module builds and tests independently of the monorepo. Publishing
the `v3.0.0` Git tag and enabling website snippets are separate release steps.

## Method index

<!-- BEGIN GENERATED METHOD INDEX -->

- [`Search`](#search) — `GET /search`
- [`Typeahead`](#typeahead) — `GET /typeahead`
- [`SearchEpisodeTitles`](#searchepisodetitles) — `GET /search_episode_titles`
- [`SpellCheck`](#spellcheck) — `GET /spellcheck`
- [`FetchRelatedSearches`](#fetchrelatedsearches) — `GET /related_searches`
- [`FetchTrendingSearches`](#fetchtrendingsearches) — `GET /trending_searches`
- [`FetchBestPodcasts`](#fetchbestpodcasts) — `GET /best_podcasts`
- [`FetchPodcastByID`](#fetchpodcastbyid) — `GET /podcasts/{id}`
- [`DeletePodcast`](#deletepodcast) — `DELETE /podcasts/{id}`
- [`FetchEpisodeByID`](#fetchepisodebyid) — `GET /episodes/{id}`
- [`BatchFetchEpisodes`](#batchfetchepisodes) — `POST /episodes`
- [`BatchFetchPodcasts`](#batchfetchpodcasts) — `POST /podcasts`
- [`FetchCuratedPodcastsListByID`](#fetchcuratedpodcastslistbyid) — `GET /curated_podcasts/{id}`
- [`FetchPodcastGenres`](#fetchpodcastgenres) — `GET /genres`
- [`FetchPodcastRegions`](#fetchpodcastregions) — `GET /regions`
- [`FetchPodcastLanguages`](#fetchpodcastlanguages) — `GET /languages`
- [`JustListen`](#justlisten) — `GET /just_listen`
- [`FetchCuratedPodcastsLists`](#fetchcuratedpodcastslists) — `GET /curated_podcasts`
- [`FetchRecommendationsForPodcast`](#fetchrecommendationsforpodcast) — `GET /podcasts/{id}/recommendations`
- [`FetchRecommendationsForEpisode`](#fetchrecommendationsforepisode) — `GET /episodes/{id}/recommendations`
- [`SubmitPodcast`](#submitpodcast) — `POST /podcasts/submit`
- [`FetchPlaylistByID`](#fetchplaylistbyid) — `GET /playlists/{id}`
- [`FetchMyPlaylists`](#fetchmyplaylists) — `GET /playlists`
- [`FetchAudienceForPodcast`](#fetchaudienceforpodcast) — `GET /podcasts/{id}/audience`
- [`FetchPodcastsByDomain`](#fetchpodcastsbydomain) — `GET /podcasts/domains/{domain_name}`
- [`CreatePlaylist`](#createplaylist) — `POST /playlists`
- [`UpdatePlaylist`](#updateplaylist) — `PUT /playlists/{id}`
- [`AddPlaylistItem`](#addplaylistitem) — `POST /playlists/{id}/items`
- [`DeletePlaylistItem`](#deleteplaylistitem) — `DELETE /playlists/{id}/items/{item_id}`
- [`UpdatePlaylistItemNotes`](#updateplaylistitemnotes) — `PUT /playlists/{id}/items/{item_id}`

<!-- END GENERATED METHOD INDEX -->

## API reference

<!-- BEGIN GENERATED API REFERENCE -->

Methods accept positional path identifiers followed by `map[string]string` parameters. Examples use the public mock; pass your API key to `NewClient` for production.

### Search

Full-text search

`GET /search`

Full-text search on episodes, podcasts, or curated lists of podcasts.
Use the `offset` parameter to paginate through search results.
The FREE plan allows to see up to 30 search results (or `offset` < 30) per query.
The PRO plan allows to see up to 300 search results (or `offset` < 300) per query.
The ENTERPRISE plan allows to see up to 10,000 search results (or `offset` < 10000) per query.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.Search(map[string]string{"q": "star wars", "sort_by_date": "0", "type": "episode", "offset": "0", "len_min": "10", "len_max": "30", "genre_ids": "68,82", "published_before": "1580172454000", "published_after": "0", "only_in": "title,description", "language": "English", "region": "", "safe_mode": "0", "unique_podcasts": "0", "interviews_only": "0", "sponsored_only": "0", "page_size": "10"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-search)

### Typeahead

Typeahead search

`GET /typeahead`

Suggest search terms, podcast genres, and podcasts.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.Typeahead(map[string]string{"q": "star wars", "show_podcasts": "1", "show_genres": "1", "safe_mode": "0"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-typeahead)

### SearchEpisodeTitles

Find individual episodes by searching for their titles

`GET /search_episode_titles`

Conduct targeted searches for individual episodes by title and refine results using the podcast id such as
Listen Notes Podcast ID, Apple Podcasts ID, Spotify ID, or RSS feed URL.
This endpoint is specially designed to streamline the import of specific episodes from platforms
like Apple Podcasts and Spotify into your application.
Compared to the GET /search endpoint, which performs full-text searches across multiple fields,
this endpoint focuses solely on episode titles for enhanced accuracy and performance.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.SearchEpisodeTitles(map[string]string{"q": "Jerusalem Demsas on The Dispossessed"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-search_episode_titles)

### SpellCheck

Spell check on a search term

`GET /spellcheck`

Suggest a list of words that correct the spelling errors of a search term. This endpoint is available only in the PRO/ENTERPRISE plan.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.SpellCheck(map[string]string{"q": "microsft stock"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-spellcheck)

### FetchRelatedSearches

Fetch related search terms

`GET /related_searches`

Suggest related search terms. The results are more comprehensive than from `GET /typeahead`. This endpoint is available only in the PRO/ENTERPRISE plan.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchRelatedSearches(map[string]string{"q": "evergrande"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-related_searches)

### FetchTrendingSearches

Fetch trending search terms

`GET /trending_searches`

Fetch up to 10 most recent trending search terms on the Listen Notes platform.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchTrendingSearches(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-trending_searches)

### FetchBestPodcasts

Fetch a list of best podcasts by genre

`GET /best_podcasts`

Get a list of curated best podcasts by genre,
which are curated by Listen Notes staffs based on various signals from the Internet, e.g.,
top charts on other podcast platforms, recommendations from mainstream media,
user activities on listennotes.com...
You can get the genre ids from `GET /genres` endpoint.
This endpoint returns same data as https://www.listennotes.com/best-podcasts/

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchBestPodcasts(map[string]string{"genre_id": "93", "page": "2", "region": "us", "sort": "listen_score", "safe_mode": "0"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-best_podcasts)

### FetchPodcastByID

Fetch detailed meta data and episodes for a podcast by id

`GET /podcasts/{id}`

Fetch detailed meta data and episodes for a specific podcast (up to 10 episodes each time).
You can use the **next_episode_pub_date** parameter to do pagination and fetch more episodes.
During pagination with **next_episode_pub_date**, an empty **episodes** array in the response signals that no more episodes are available.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchPodcastByID("4d3fe717742d4963a85562e9f84d8c79", map[string]string{"next_episode_pub_date": "1479154463000", "sort": "recent_first"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-podcasts-id)

### DeletePodcast

Request to delete a podcast

`DELETE /podcasts/{id}`

Podcast hosting services can use this endpoint to streamline the process of podcast deletion on behave of their users (podcasters). We will review the deletion request within 12 hours. If the podcast is already deleted, the "status" field in the response will be "deleted". Otherwise, the status field will be "in review". If you want to get a notification once the podcast is deleted, you can configure a webhook url in the dashboard: listennotes.com/api/dashboard/#webhooks

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.DeletePodcast("4d3fe717742d4963a85562e9f84d8c79", map[string]string{"reason": "the podcaster wants to delete it"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#delete-api-v2-podcasts-id)

### FetchEpisodeByID

Fetch detailed meta data for an episode by id

`GET /episodes/{id}`

Fetch detailed meta data for a specific episode.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchEpisodeByID("6b6d65930c5a4f71b254465871fed370", map[string]string{"show_transcript": "1"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-episodes-id)

### BatchFetchEpisodes

Batch fetch basic meta data for episodes

`POST /episodes`

Batch fetch basic meta data for up to 10 episodes. This endpoint could be used to implement custom playlists for individual episodes. For detailed meta data of an individual episode, you need to use `GET /episodes/{id}`. This endpoint is available only in the PRO/ENTERPRISE plan.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.BatchFetchEpisodes(map[string]string{"ids": "c577d55b2b2b483c969fae3ceb58e362,0f34a9099579490993eec9e8c8cebb82"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#post-api-v2-episodes)

### BatchFetchPodcasts

Batch fetch basic meta data for podcasts

`POST /podcasts`

Batch fetch basic meta data for up to 10 podcasts.
This endpoint could be used to build something like OPML import,
allowing users to import a bunch of podcasts via rss urls.
For detailed meta data (including episodes) of an individual podcast, you need to use `GET /podcasts/{id}`. This endpoint is available only in the PRO/ENTERPRISE plan.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.BatchFetchPodcasts(map[string]string{"ids": "3302bc71139541baa46ecb27dbf6071a,68faf62be97149c280ebcc25178aa731,37589a3e121e40debe4cef3d9638932a,9cf19c590ff0484d97b18b329fed0c6a", "rsses": "https://rss.art19.com/recode-decode,https://rss.art19.com/the-daily,https://www.npr.org/rss/podcast.php?id=510331,https://www.npr.org/rss/podcast.php?id=510331", "itunes_ids": "1457514703,1386234384,659155419", "spotify_ids": "3DDfEsKDIDrTlnPOiG4ZF4,4qDNe5Gvl1XxdLinUGEXrC,23NZCM4ik6o3UYkM473Itz", "show_latest_episodes": "1", "next_episode_pub_date": "1557394247000"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#post-api-v2-podcasts)

### FetchCuratedPodcastsListByID

Fetch a curated list of podcasts by id

`GET /curated_podcasts/{id}`

Get detailed meta data of all podcasts in a specific curated list.
This endpoint returns same data as https://www.listennotes.com/curated-podcasts/

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchCuratedPodcastsListByID("SDFKduyJ47r", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-curated_podcasts-id)

### FetchPodcastGenres

Fetch a list of podcast genres

`GET /genres`

Get a list of podcast genres that are supported in Listen Notes.
The genre id can be passed to other endpoints as a parameter to get podcasts in a specific genre,
e.g., `GET /best_podcasts`, `GET /search`...
You may want to cache the list of genres on the client side.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchPodcastGenres(map[string]string{"top_level_only": "1"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-genres)

### FetchPodcastRegions

Fetch a list of supported countries/regions for best podcasts

`GET /regions`

It returns a dictionary of country codes (e.g., us, gb...) & country names (United States, United Kingdom...). The country code is used in the query parameter **region** of `GET /best_podcasts`.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchPodcastRegions(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-regions)

### FetchPodcastLanguages

Fetch a list of supported languages for podcasts

`GET /languages`

Get a list of languages that are supported in Listen Notes database. You can use the language string as query parameter in `GET /search`.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchPodcastLanguages(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-languages)

### JustListen

Fetch a random podcast episode

`GET /just_listen`

Recently published episodes are more likely to be fetched. Good luck!

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.JustListen(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-just_listen)

### FetchCuratedPodcastsLists

Fetch curated lists of podcasts

`GET /curated_podcasts`

A bunch of curated lists from online media. For each list, you'll get basic info of up to 5 podcasts. To get detailed meta data of all podcasts in a specific list, you need to use `GET /curated_podcasts/{id}`. We add new curated lists to the database on a daily basis.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchCuratedPodcastsLists(map[string]string{"page": "2"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-curated_podcasts)

### FetchRecommendationsForPodcast

Fetch recommendations for a podcast

`GET /podcasts/{id}/recommendations`

Fetch up to 8 podcast recommendations based on the given podcast id.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchRecommendationsForPodcast("25212ac3c53240a880dd5032e547047b", map[string]string{"safe_mode": "0"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-podcasts-id-recommendations)

### FetchRecommendationsForEpisode

Fetch recommendations for an episode

`GET /episodes/{id}/recommendations`

Fetch up to 8 episode recommendations based on the given episode id.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchRecommendationsForEpisode("254444fa6cf64a43a95292a70eb6869b", map[string]string{"safe_mode": "0"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-episodes-id-recommendations)

### SubmitPodcast

Submit a podcast to Listen Notes database

`POST /podcasts/submit`

Podcast hosting services can use this endpoint to help your users directly submit a new podcast to Listen Notes database. If the podcast doesn't exist in the database, "status" in the response will be "in review", and we'll review it within 12 hours. If the podcast exists, "status" in the response will be "found". If this submission is rejected, "status" in the response will be "rejected". You can use `POST /podcasts` to check if multiple podcasts exist in the database. If you want to get a notification once the podcast is accepted, you can either specify the "email" parameter or configure a webhook url in the dashboard: listennotes.com/api/dashboard/#webhooks

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.SubmitPodcast(map[string]string{"rss": "https://feeds.megaphone.fm/committed", "email": "hello@example.com"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#post-api-v2-podcasts-submit)

### FetchPlaylistByID

Fetch a playlist's info and items (i.e., episodes or podcasts).

`GET /playlists/{id}`

A playlist can contain both episodes and podcasts, shown in separate views,
just like playlists created via listennotes.com/listen/.
This endpoint fetches items from the saved default view unless **type** is specified.
The response type and listennotes_url describe the selected view.
You can use the **last_pub_date_ms** parameter to do pagination and fetch more items.
A playlist can be **public** (discoverable on ListenNotes.com),
**unlisted** (accessible to anyone who knows the playlist id),
or **private** (accessible when the API admin has active playlist membership).
Public and unlisted playlists can also be fetched by ID regardless of their owner.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchPlaylistByID("m1pe7z60bsw", map[string]string{"type": "episode_list", "last_timestamp_ms": "0", "sort": "recent_added_first"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-playlists-id)

### FetchMyPlaylists

Fetch a list of your playlists.

`GET /playlists`

This endpoint lists playlists with an active membership for the API admin, including playlists they created or joined.
Each playlist includes its saved default **type** and a **listennotes_url** for that view.
You can use the **page** parameter to do pagination and fetch more playlists.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchMyPlaylists(map[string]string{"sort": "recent_added_first", "page": "1"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-playlists)

### FetchAudienceForPodcast

Fetch audience demographics for a podcast

`GET /podcasts/{id}/audience`

Fetch audience demographics for a podcast - 1) directly measured on the Listen Notes platform; 2) only supports audience breakdown by regions for now; 3) not every podcast has data.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchAudienceForPodcast("25212ac3c53240a880dd5032e547047b", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-podcasts-id-audience)

### FetchPodcastsByDomain

Fetch podcasts by a publisher's domain name

`GET /podcasts/domains/{domain_name}`

Fetch podcasts by a publisher's domain name, e.g., nytimes.com, wondery.com, npr.org...
Each request will return up to 10 podcasts. You can use the `page` parameter to paginate.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.FetchPodcastsByDomain("nytimes.com", map[string]string{"page": "1"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#get-api-v2-podcasts-domains-domain_name)

### CreatePlaylist

Create a playlist.

`POST /playlists`

Create an empty playlist owned by the API admin. Name is required; description defaults to an empty string, visibility defaults to public, and type defaults to episode_list. Set type to podcast_list to make podcasts the default view. The response includes the saved type and its listennotes_url.

Only playlists owned by your admin API account can be modified; contributor membership does not grant write access.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.CreatePlaylist(map[string]string{"name": "My favorite podcasts", "description": "Podcasts and episodes to revisit.", "visibility": "public", "type": "episode_list"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#post-api-v2-playlists)

### UpdatePlaylist

Update playlist metadata.

`PUT /playlists/{id}`

Update any subset of name, description, visibility, and type. Omitted fields remain unchanged; at least one field is required. Switching to private rotates the playlist RSS secret. Type selects the saved default view (episode_list or podcast_list) and the returned listennotes_url; changing it preserves all existing episodes and podcasts.

Only playlists owned by your admin API account can be modified; contributor membership does not grant write access.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.UpdatePlaylist("m1pe7z60bsw", map[string]string{"name": "My favorite podcasts", "description": "Podcasts and episodes to revisit.", "visibility": "public", "type": "podcast_list"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#put-api-v2-playlists-id)

### AddPlaylistItem

Add an episode or podcast to a playlist.

`POST /playlists/{id}/items`

Provide exactly one non-empty episode_id or podcast_id; an empty unused ID field is ignored. Invalid ID formats return 400 and identify the field. A missing episode or podcast returns 404 with an error such as "Episode not found: {episode_id}." or "Podcast not found: {podcast_id}.". Existing active items are reused (200); new or restored items return 201. Omitted notes preserve existing notes, including when restoring a deleted item; supplied notes replace them.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.AddPlaylistItem("m1pe7z60bsw", map[string]string{"episode_id": "e53e6992a5b7492f9ea6fcd85d9ad95f", "notes": "Worth a listen."})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#post-api-v2-playlists-id-items)

### DeletePlaylistItem

Remove an item from a playlist.

`DELETE /playlists/{id}/items/{item_id}`

Delete a playlist item. Repeating deletion of the same item succeeds. This does not delete the episode or podcast from the podcast database.

Only playlists owned by your admin API account can be modified; contributor membership does not grant write access.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.DeletePlaylistItem("m1pe7z60bsw", "23", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#delete-api-v2-playlists-id-items-item_id)

### UpdatePlaylistItemNotes

Update notes for a playlist item.

`PUT /playlists/{id}/items/{item_id}`

Replace item notes, or send an empty string to clear them. The item ID and added_at_ms remain unchanged.

Only playlists owned by your admin API account can be modified; contributor membership does not grant write access.

```go
package main

import (
	"fmt"
	listennotes "github.com/ListenNotes/podcast-api-go/v3"
)

func main() {
	client := listennotes.NewClient("")
	response, err := client.UpdatePlaylistItemNotes("m1pe7z60bsw", "23", map[string]string{"notes": ""})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.ToJSON())
}
```

[Full API documentation](https://www.listennotes.com/api/docs/#put-api-v2-playlists-id-items-item_id)

<!-- END GENERATED API REFERENCE -->
