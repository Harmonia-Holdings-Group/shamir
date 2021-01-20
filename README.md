# shamir
Implementation of Shamir's secret sharing scheme
# Overview
### Encryption
1. Based on an user-provided master password, a key **K** is generated using **SHA-256**.
2. Content is encrypted using **AES-256** with **K** as key.
3. A **t - 1** degree polynomial is randomly generated that will later be used to generate **n** key shares by randomly taking **xi** and its evaluation against the polynomial **P(xi)**. Key shares are pairs **(xi, P(xi))** that can be used to recover **K** by evaluating **P(0)** using _Horner's Method_.

### Decryption
With **t** of the **n** key shares generated, the **K** key used to encrypt the content, can be recovered. 

1. Using _Horner's Method_ of polynomial evaluation, we evaluate **P(0)**, thus recovering **K**.
2. Using **K** we can decrypt the content.
