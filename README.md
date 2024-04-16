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