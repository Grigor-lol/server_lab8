package database

const (
	AddPlatformQuery = `INSERT INTO video_games.platform (platform_name)
						  VALUES ((?));`

	DeletePlatformQuery1 = `
DELETE FROM video_games.region_sales WHERE game_platform_id IN 
                                           (SELECT id FROM video_games.game_platform WHERE platform_id IN
                                                                                           (SELECT id FROM video_games.platform WHERE platform_name =?));
		`
	DeletePlatformQuery2 = `DELETE FROM video_games.game_platform WHERE platform_id IN
						                                              (SELECT id FROM video_games.platform WHERE platform_name =?);
			`
	DeletePlatformQuery3 = `DELETE FROM video_games.platform WHERE platform_name = ?;`

	AddGameQuery = `INSERT INTO video_games.game (genre_id, game_name) 
							VALUES ((SELECT id from video_games.genre where genre_name = ?), ?)`

	UpdateGameReleaseYear = `UPDATE video_games.game_platform
							SET release_year = ?
							WHERE game_publisher_id = (SELECT id from video_games.game_publisher where game_id IN (SELECT id FROM video_games.game where game_name = ?)) 						                                                         
							        `
)
