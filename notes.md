# Linode Setup

1. Make a linode account, deploy a debian image
1. SSH into the debian image
1. Just copy this `.bash_profile` to `~/.bash_profile`, to make your life easier:

        alias ls='ls -a --color=auto'
        source ~/.gvm/scripts/gvm
        gvm use go1.10.1
        export PATH=$PATH:`go env GOPATH`/bin

1. Run 

        apt-get update
        apt-get install git
        apt-get install unzip
        apt-get install golang

1. But apt-get usually doesn't have the most up-to-date version of golang! So why did we intall it? You'll see...
1. Install gvm (go version manager) 

        bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)

1. Now exit your ssh session, and log back in (so that .bash_profile is `source`'d)
1. Try running `gvm version` to make sure it's all good.
1. The reason we installed an older version of golang, is that golang itself is *written* in golang after 1.4, so we need a post-1.4 version of golang to build it from source.
1. Set our gvm bootstrap env var, so that we can build a modern version of Go, with `export GOROOT_BOOTSTRAP=$GOROOT`
1. Install the new version of golang `gvm install go<version_u_like>`
1. Now that you have a new version of go installed, we can remove the old one with 

        apt-get remove golang && apt-get autoremove

1. Ok ok ok now we have go installed.  We still need protobufs though.  First thing's first, let's make our lives easier by making our gopath accessible `ln -s ~/go $GOPATH`.
1. Grab a fresh copy of protoc from https://github.com/google/protobuf/releases.  Find a link that matches your linode (probably `linux-x86_64`), and run `wget <link>` on the link for the release.
1. Go get the protoc -ut for golang `go get -u github.com/golang/protobuf/protoc-gen-go`.  At this point, you might need to exit and log back in (to source our `bash_profile`).  Alternatively just `source ~/.bash_profile`.
1. At this point, we have basically everything we need!! Woah!  Now we just need to get our github ssh keys all sync'd.  Follow the guide <a href="https://help.github.com/articles/adding-a-new-ssh-key-to-your-github-account/">here</a>.  Then you can `git clone` whatever private repos you have access to.