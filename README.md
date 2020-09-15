# go-acapy-client

A library for interacting with ACA-py in Go.

## Implemented Endpoints

### Connections

| Function Name                  | Method                               | Endpoint                                     | Implemented |
|----------------------------------------------|------------------------------|------------------------------|------------------------------|
| QueryConnections            | GET                              | /connections                                 | :heavy_check_mark: |
| GetConnection | GET | /connections/{id} | :heavy_check_mark: |
| CreateInvitation | POST           | /connections/create-invitation               | :heavy_check_mark: |
| ReceiveInvitation | POST          | /connections/receive-invitation              | :heavy_check_mark: |
| AcceptInvitation | POST      | /connections/{id}/accept-invitation          | :heavy_check_mark: |
| AcceptRequest | POST         | /connections/{id}/accept-request             | :heavy_check_mark: |
| RemoveConnection    | POST \<why though :man_facepalming:> | /connections/{id}/remove                     | :heavy_check_mark: |
| SendBasicMessage    | POST                | /connections/send-message                    | :heavy_check_mark: |
| SendPing               | POST               | /connections/send-ping                       | :heavy_check_mark: |
