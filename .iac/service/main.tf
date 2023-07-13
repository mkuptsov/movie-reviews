data "moviereviews_user" "johndoe" {
    username = "johndoe" 
}

resource "moviereviews_user_role" "johndoe" {
    user_id = data.moviereviews_user.johndoe.id
    role = "admin"
}

