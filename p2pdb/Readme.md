# Peer to Peer Node database

Database is meant to store data with every record being signed 2 times to be considered valid(Can be signed more.. Will come later.)  
Black ball public keys... All keys you come in contact with are valid to decrypt messages... (All connections / data transfer is done via encryted messaging.)

All data is stored via encrypted records.. some fields un encrypted others encrypted... Might use Table to correlate validity of store.

Example UserA sends UserB message "Hello" The message is sent And UserB decrypts it. UserB store that + pub key in the in db With him denoted as signer.
UserB sends that message to other interested parties.. Eg back to UserA. UserA sees UserB signed it As a signer verfiying UserA(His message) and UserA
stores his original message in its db with him as author 

# Web of trust auth.


                                                    +-------------------------------+
                                                    |                               |
                                                    |                               |
                                                    |                               |                                       +-------------------------------+
                                                    |                               |                                       |                               |
+-------------------------------+                   |         UserA:PubKey3         |                                       |                               |
|                               |                   |                               |                                       |      UserC:PubKey8            |
|       Assumptions:            |                   |                               |                                       |                               |
|       Nodes always communicate|                   |                               |                                       |                               |
|        to one another directly.                   |                               |                                       |                               |
|        A-C C-A                |                   |                               |                                       |                               |
|       C-B B-C never           |                   |                               |                                       |                               |
|       C-B through A           |                   |                               |                                       |                               |
|                               |                   |                               |                                       |                               |
|       Nodes all can store     |                   |                               |                                       |                               |
|       data for antother node  |                   |                               |                                       |                               |
|       Data is only ever guarnteed                 |                               |                                       |                               |
|       to be on at least n nodes                   |                               |                                       |                               |
|       based on policy set for message             |                               |                                       |                               |
|       Ex: I create messageA it is only valid      |                               |                                       |                               |
|       after n people have signed                  |                               |                                       +-------------------------------+
|       Ex "Hello" 3Sigs becomes valid after        |                               |
|       3 times being signed.   |                   |                               |
|       Which means it will be stored on            +-------------------------------+
|       at least 3 nodes to exist.
|                               |
|       But 1 sign only requires me to send it.                                                                           +-----------------------------+
|       Records are stored in a |                                                                                         |                             |
|       db that only stores its ID                    +------------------------------+                                    |                             |
|        and raw content.       |                     |                              |                                    |     UserD:PubKey19          |
|        Other areas have stored actual               |                              |                                    |                             |
|        signed message n times up.                   |    UserB:PubKey6             |                                    |                             |
|                               |                     |                              |                                    |                             |
|                               |                     |                              |                                    |                             |
|                               |                     |                              |                                    |                             |
|                               |                     |                              |                                    |                             |
|                               |                     |                              |                                    |                             |
|                               |                     |                              |                                    |                             |
+-------------------------------+                     |                              |                                    |                             |
                                                      |                              |                                    |                             |
                                                      |                              |                                    |                             |
                                                      |                              |                                    |                             |
                                                      |                              |                                    |                             |
                                                      |                              |                                    +-----------------------------+
                                                      |                              |
                                                      +------------------------------+




