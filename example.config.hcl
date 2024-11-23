application {
  check_for_updates = true
}

output {
  dir       = "downloads"
  clean_dir = false
  format    = "auto" # auto | jpeg | png | avif | webp
}

download {
  preload_next_chapters = 2
  page_batch_size = 4
}

site {
  "shonenjumpplus.com" {
    cookie_string = ""
  }

  "comic-zenon.com" {
    cookie_string = ""
  }

  "pocket.shonenmagazine.com" {
    cookie_string = ""
  }

  "comic-gardo.com" {
    cookie_string = ""
  }

  "magcomi.com" {
    cookie_string = ""
  }

  "comic-action.com" {
    cookie_string = ""
  }

  "comic-days.com" {
    cookie_string = ""
  }

  "kuragebunch.com" {
    cookie_string = ""
  }

  "viewer.heros-web.com" {
    cookie_string = ""
  }

  "www.cmoa.jp" {
    cookie_string = ""
  }

  "www.corocoro.jp" {
    cookie_string = ""
  }
}
