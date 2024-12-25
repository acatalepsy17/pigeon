package routes

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Endpoint struct {
	DB *gorm.DB
}

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	endpoint := Endpoint{DB: db}
	api := app.Group("/api/v1")

	// authentication
	authRouter := api.Group("/auth")
	authRouter.Post("/register", endpoint.Register)
	authRouter.Post("/verify-email", endpoint.VerifyEmail)
	authRouter.Post("/resend-verification-email", endpoint.ResendVerificationEmail)
	authRouter.Post("/send-password-reset-otp", endpoint.SendPasswordResetOtp)
	authRouter.Post("/set-new-password", endpoint.SetNewPassword)
	authRouter.Post("/login", endpoint.Login)
	authRouter.Post("/refresh", endpoint.Refresh)
	authRouter.Get("/logout", endpoint.Logout)

	// user profile
	profilesRouter := api.Group("/profiles", endpoint.AuthMiddleware)
	profilesRouter.Get("/cities", endpoint.RetrieveCities)
	profilesRouter.Get("", endpoint.GuestMiddleware, endpoint.RetrieveUsers)
	profilesRouter.Get("/profile/:username", endpoint.RetrieveUserProfile)
	profilesRouter.Patch("/profile", endpoint.UpdateProfile)
	profilesRouter.Post("/profile", endpoint.DeleteUser)
	profilesRouter.Get("/friends", endpoint.RetrieveFriends)
	profilesRouter.Get("/friends/requests", endpoint.RetrieveFriendRequests)
	profilesRouter.Post("/friends/requests", endpoint.SendOrDeleteFriendRequest)
	profilesRouter.Put("/friends/requests", endpoint.AcceptOrRejectFriendRequest)
	profilesRouter.Get("/notifications", endpoint.RetrieveUserNotifications)
	profilesRouter.Post("/notifications", endpoint.ReadNotification)

	// newsfeed
	feedRouter := api.Group("/feed", endpoint.AuthMiddleware)
	feedRouter.Get("/posts", endpoint.RetrievePosts)
	feedRouter.Post("/posts", endpoint.CreatePost)
	feedRouter.Get("/posts/:slug", endpoint.RetrievePost)
	feedRouter.Put("/posts/:slug", endpoint.UpdatePost)
	feedRouter.Delete("/posts/:slug", endpoint.DeletePost)
	feedRouter.Get("/reactions/:focus/:slug", endpoint.RetrieveReactions)
	feedRouter.Post("/reactions/:focus/:slug", endpoint.CreateReaction)
	feedRouter.Delete("/reactions/:id", endpoint.DeleteReaction)
	feedRouter.Get("/posts/:slug/comments", endpoint.RetrieveComments)
	feedRouter.Post("/posts/:slug/comments", endpoint.CreateComment)
	feedRouter.Get("/comments/:slug", endpoint.RetrieveCommentWithReplies)
	feedRouter.Post("/comments/:slug", endpoint.CreateReply)
	feedRouter.Put("/comments/:slug", endpoint.UpdateComment)
	feedRouter.Delete("/comments/:slug", endpoint.DeleteComment)
	feedRouter.Get("/replies/:slug", endpoint.RetrieveReply)
	feedRouter.Put("/replies/:slug", endpoint.UpdateReply)
	feedRouter.Delete("/replies/:slug", endpoint.DeleteReply)

	// communication
	chatRouter := api.Group("/chats", endpoint.AuthMiddleware)
	chatRouter.Get("", endpoint.RetrieveUserChats)
	chatRouter.Post("", endpoint.SendMessage)
	chatRouter.Get("/:chat_id", endpoint.RetrieveMessages)
	chatRouter.Patch("/:chat_id", endpoint.UpdateGroupChat)
	chatRouter.Delete("/:chat_id", endpoint.DeleteGroupChat)
	chatRouter.Put("/messages/:message_id", endpoint.UpdateMessage)
	chatRouter.Delete("/messages/:message_id", endpoint.DeleteMessage)
	chatRouter.Post("/groups/group", endpoint.CreateGroupChat)

	// websocket
	api.Get("/ws/notifications", websocket.New(endpoint.NotificationSocket))
	api.Get("/ws/chats/:id", websocket.New(endpoint.ChatSocket))
}
