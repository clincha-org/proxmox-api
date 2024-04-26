# proxmox-api

I need to make an API to connect to Proxmox so that I can make a Terraform Provider to interact with it. Everything is written in Go(lang) so I need to learn that as well. I've checked a few videos on YouTube and looked at some documentation, it should be pretty straightforward but a lot of work. I want to do as much testing as I can so that the code is really robust. I've created some Vagrant boxes for that in another repository.

While making the Terraform Provider I got to [this tutorial](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider-configure) before I realised that I needed to make the API as well.

A useful search term for the work I'm doing in this repo is "go REST client". Looks like the http library in Go has had a recent update, so I'm finding a few tutorials are out of date. 

The Terraform tutorial that I'm following wants me to call this function called "NewClient" to initiate the connection. That will need to be created at some point but I think I just need to use the `main.go` file to run everything right now.

---

I've created a `main_test.go` where I can run tests from instead of having a `main.go` to run from each time. I've also implemented a `pkg` folder after reading about the [project layout](https://github.com/golang-standards/project-layout) that Go projects should have. Then I've added a `proxmox` folder to that and a file called `client.go` which is going to deal with the `NewClient` function that Terraform calls. I'm going to implement the connection logic there and then hook it back up to the Terraform provider code and see if I can get it working end to end before I go and make anything else. That way I can make a test in Terraform that does a full authentication test from the provider through this client API and finally into the actual Proxmox server.

---

I've hit an issue where I keep getting `401 No ticket` as a response. I think this is to do with the authentication headers that I'm sending. This [forum post](https://forum.proxmox.com/threads/working-with-api-getting-401-no-ticket.75108/) looks like it has the answer. 

It wasn't the issue... the issue was that I had a `/` at the end of the request. So instead of `access/ticket` it was `access/ticket/`. I hate computers sometimes.

Got some basic tests in place for authentication. Not sure that I like the current implementation, but we can grow from here. I think the next stage is going to be connecting everything back up to the terraform provider. I also want to get the vagrant machines built into the tests so that I don't need to spin them up separately. A GitHub workflow that runs the tests on commit would be another good step.


---

I have connected the Proxmox provider up to this API and done a dummy node datasource. I know that it can log in and connect everything together successfully. It's nice to see the terraform output as green.

```text
/usr/bin/terraform plan
╷
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI
│ configuration:
│  - hashicorp.com/edu/proxmox in /home/clincha/go/bin
│ 
│ The behavior may therefore not match any released version of the provider
│ and applying changes may cause the state to become incompatible with
│ published releases.
╵
data.proxmox_node.example: Reading...
data.proxmox_node.example: Read complete after 0s

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration
and found no differences, so no changes are needed.

Process finished with exit code 0
```

---

I am now implementing a data source. I've added some of the code to the provider but I now need to implement a "GetNodes" function for it to call. 

I've added the code now. IntelliJ has a cool thing that when I paste json into a go file it converts it into a struct for me. I've used that to generate the structs for me which is way nicer.

I've removed the checks I was doing for the capabilities of the node because it was hard to do reflection and figure out the names and values of the capabilities struct. I can revisit this at some point, but I want to continue. I think the tests will become more robust when I'm able to layer them on top of each-other. One idea I have is for a test that creates two nodes, puts them in a cluster and then calls getNodes and makes sure there are two of them.

---

I've implemented the data source for nodes! I get the data back from the API and it's coming out as Terraform output!

```text
Changes to Outputs:
  + edu_nodes = {
      + node = {
          + data = [
              + {
                  + cpu             = 0.00597014925373134
                  + disk            = 3672838144
                  + id              = "node/pve"
                  + level           = ""
                  + max_cpu         = 1
                  + max_disk        = 19136331776
                  + max_memory      = 2043310080
                  + memory          = 1098805248
                  + node            = "pve"
                  + ssl_fingerprint = "62:49:13:7F:E8:01:D7:3F:09:68:5C:A5:F3:95:02:39:8A:72:42:40:4A:20:50:9F:53:45:B5:CF:F8:0A:F7:B1"
                  + status          = "online"
                  + type            = "node"
                  + uptime          = 4294
                },
            ]
        }
    }
```

How cool is that!? 

Anyway, issues with it so far
- The data wrapper is pretty rubbish
- Not sure what those numbers all mean
- I could probably make it better with some descriptions

All fixable issues but let's get it done before moving on to the next stage. The issues I want to fix more generally are:
- The terraform code isn't in a repository that I control
- There aren't any tests in the Terraform code for me to run
- This is all happening in my main PC so I can't get it working on my laptop
- This is all working against Proxmox VE 7 but not 8
- The rest of the API needs to be implemented!

---

I've sorted out the data wrapper issue! This is much nicer output:

```text
Changes to Outputs:
  + edu_nodes = {
      + nodes = [
          + {
              + cpu             = 0.0029940119760479
              + disk            = 3673042944
              + id              = "node/pve"
              + level           = ""
              + max_cpu         = 1
              + max_disk        = 19136331776
              + max_memory      = 2043310080
              + memory          = 1092751360
              + node            = "pve"
              + ssl_fingerprint = "62:49:13:7F:E8:01:D7:3F:09:68:5C:A5:F3:95:02:39:8A:72:42:40:4A:20:50:9F:53:45:B5:CF:F8:0A:F7:B1"
              + status          = "online"
              + type            = "node"
              + uptime          = 8005
            },
        ]
    }
```

---

The next part of the tutorial is about making resources. I've had a look through the API documentation for Proxmox and it looks like the most simple resource to create is probably a network. I'm going to try and create one of those. I have started the Terraform Provider side but I can't go much further without having a dig into the actual API and getting some code working to pull that data through.

