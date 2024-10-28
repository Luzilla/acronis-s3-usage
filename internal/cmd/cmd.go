package cmd

import "github.com/urfave/cli/v2"

func BucketCommand() *cli.Command {
	return &cli.Command{
		Name:    "buckets",
		Aliases: []string{"b"},
		Usage:   "list buckets",
		Action:  listBuckets,
		Flags: []cli.Flag{
			emailFlag(),
		},
		Subcommands: []*cli.Command{
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Action:  deleteBucket,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "bucket",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "confirm",
						Value: false,
					},
				},
			},
			{
				Name:    "show",
				Aliases: []string{"s"},
				Action:  showBucket,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "bucket",
						Required: true,
					},
				},
			},
		},
	}
}

func StatsCommand() *cli.Command {
	return &cli.Command{
		Name:    "stats",
		Aliases: []string{"s"},
		Usage:   "list stats",
		Flags:   []cli.Flag{
			// &cli.TimestampFlag{
			// 	Name:     "from",
			// 	Layout:   "2006-01-02",
			// 	Timezone: time.UTC,
			// 	Required: false,
			// },
		},
		Action: showStats,
	}
}

func UsersCommand() *cli.Command {
	return &cli.Command{
		Name:    "users",
		Aliases: []string{"u"},
		Usage:   "manage users",
		Action:  users,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "usage",
				Value: false,
			},
		},
		Subcommands: []*cli.Command{
			{
				Name: "create",
				Flags: []cli.Flag{
					emailFlag(),
				},
				Action: createUser,
			},
			{
				Name: "delete",
				Flags: []cli.Flag{
					emailFlag(),
				},
				Action: deleteUser,
			},
			{
				Name: "lock",
				Flags: []cli.Flag{
					emailFlag(),
				},
				Action: lockUser,
			},
			{
				Name: "unlock",
				Flags: []cli.Flag{
					emailFlag(),
				},
				Action: unlockUser,
			},
			{
				Name: "show",
				Flags: []cli.Flag{
					emailFlag(),
				},
				Action: showUser,
			},
			{
				Name: "create-key",
				Flags: []cli.Flag{
					emailFlag(),
				},
				Action: createKey,
			},
			{
				Name: "revoke-key",
				Flags: []cli.Flag{
					emailFlag(),
					keyIDFlag(),
				},
				Action: revokeKey,
			},
			{
				Name: "rotate-key",
				Flags: []cli.Flag{
					emailFlag(),
					keyIDFlag(),
				},
				Action: rotateKey,
			},
			{
				Name: "limits",
				Flags: []cli.Flag{
					emailFlag(),
				},
				Action: userLimits,
			},
		},
	}
}
