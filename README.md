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

---

This was tricky, but I learned a lot that will help me out with debugging and make my code way better at reporting errors.

First thing is that HTTP responses normally return a body of some description, even if the response is not a 200. I spent a long time trying to figure out what the issue was when I only got a 400 response. When I looked into the body of the response it gave me the exact error. There was some issue with the JSON that I was passing in and I realised that I hadn't set the Content-Type header that it wanted.

Having all the tests is also really useful. Being able to run everything as soon as I make a change is awesome.

Anyway, the network now has a create and delete method. I've only tried creating and deleting a bridge, but I'm going to do the Terraform side now and see how we get on.

---

I found a linter. I want to do this when I set up the GitHub actions to trigger the builds for this. https://golangci-lint.run/

---

I started to work on the Terraform side of things and then realised that I needed to return the Network object when I created it. So I added a function that gets a specific interface given the interface name. I can then call that at the end of the create which has the added benefit that it checks with the API to make sure that the network does actually exist.

---

I have most of the Terraform code in place now, and I'm running the terraform apply command. However, I've hit my first error with pointers! This will be a fun one to debug. I'm sure I'll learn a lot, but I think this one could be painful...

```text
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0xa85687]
```

---

I watched a YouTube video about pointers and the nil pointer dereference error and learned a lot. It was great to see a basic example of how objects in Go work, so I could apply it to this. I didn't even realise that without a pointer Go copies the object to the function. I would have expected it to modify the object I send it, but it creates a new one and modifies that which is rarely useful.

I have got the apply working now :)

---

I have got the update function working now and I have also created a reload function to reload the interface everytime that it has changed. Next step is to implement all of that over to the Terraform provider. I have decided that the first release is going to be the network work. Before that I have a lot of work to do:

- Make the error messages clearer and consistent
- Add logging
- Test other interface types
- Ensure the update method works as expected without needing to remove a bunch of fields
- Make the Terraform network resource work
- Lint the code
- Create a build and release workflow in GitHub
- Decide on the format of this diary/log
- Add documentation to the functions
- Add a CONTRIBUTING.md file to the project
- Start a new repository for the Terraform code so that I don't have all the HashiCups example code
- Move this list into a set of GitHub issues and create a milestone
- Have the testing code create the vagrant boxes?