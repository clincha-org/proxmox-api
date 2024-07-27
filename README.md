# proxmox-api

## Mantras

The API must work like Terraform. If a value is not set, the zero value is assumed.

## Attribute Flow

For each attribute, we must think about three flows.

1. The attribute is set to **some value** and sent to the API. 
2. The attribute is set to a **zero value** and sent to the API. 
3. The attribute is not set (`nil`) and sent to the API.

For each of these flows there is an associated read from the API. We need to have a clear understanding of what the API should return in each of these cases. The API should be designed to be consumed by Terraform. 

For this example I will use the `comment` attribute of the `Network` object.

1. The `comment` attribute is set to "test" in the `Network` object and sent to the API.
2. The API will return the `Network` object with the `comment` attribute set to "test".
3. The `comment` attribute is set to "" in the `Network` object and sent to the API.
4. The API will return the `Network` object with the `comment` attribute set to "".
5. The `comment` attribute is not set (`nil`) in the `Network` object and sent to the API.
6. The API will return the `Network` object with the `comment` attribute set to "".

This example assumes that we are using a pointer for the comment attribute. If we do not use a pointer, then we do not need to worry about the `nil` value. 

```text
"test" --> create --> "test"
"test" <-- read <-- "test"

"" --> update --> ""
"" <-- read <-- ""

nil --> update --> ""
"" <-- read <-- ""
```
