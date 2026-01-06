# blog_aggregator

A multi-player command line tool for aggregating RSS feeds and viewing the posts.

## Installation

Make sure you have the latest [Go toolchain](https://golang.org/dl/) installed as well as a local Postgres database. You can then install `blog_aggregator` with:

```bash
go install ...
```

## Config

Create a `.gatorconfig.json` file in your home directory with the following structure:

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the values with your database connection string.

## Usage

Create a new user:

```bash
blog_aggregator register <name>
```

Add a feed:

```bash
blog_aggregator addfeed <url>
```

Start the aggreblog_aggregator:

```bash
blog_aggregator agg 30s
```

View the posts:

```bash
blog_aggregator browse [limit]
```

There are a few other commands you'll need as well:

- `blog_aggregator login <name>` - Log in as a user that already exists
- `blog_aggregator users` - List all users
- `blog_aggregator feeds` - List all feeds
- `blog_aggregator follow <url>` - Follow a feed that already exists in the database
- `blog_aggregator unfollow <url>` - Unfollow a feed that already exists in the database
