terraform {
    required_providers {
        moviereviews = {
            source = "localhost/providers/moviereviews"
        }
    }
}

provider "moviereviews" {
}
