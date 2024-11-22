application {
  check_for_updates = true
}

output {
  dir       = "downloads"
  clean_dir = false
  format    = "auto" # auto | jpeg | png | avif | webp
}

download {
  concurrent_processes = 4
}

site {
  "comic-walker.com" {}

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

  "tonarinoyj.jp" {}

  "comic-ogyaaa.com" {}

  "comic-action.com" {
    cookie_string = ""
  }

  "comic-days.com" {
    cookie_string = ""
  }

  "comic-growl.com" {}

  "comic-earthstar.com" {}

  "comicborder.com" {}

  "comic-trail.com" {}

  "kuragebunch.com" {
    cookie_string = ""
  }

  "viewer.heros-web.com" {
    cookie_string = ""
  }

  "www.sunday-webry.com" {}

  "www.cmoa.jp" {
    cookie_string = ""
  }
}
