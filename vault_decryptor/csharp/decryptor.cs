// Initialize MetaMaskCrypto
MetaCrypto crypto = new();

// The password used for decryption
string password = "Admin231@";
// The payload or data to be decrypted
string payload = "{\"data\":\"0UTpH53VzgKxrOF3aFBYyByZNC6Xy7Xp8A3pWd04GnW7ei8s6bVgkUNNSwnTei96pCxGd3CH26CBnqmuXj4QLd/rEyUocpplR6sEXYCpt+QvS951e8QwsgzWNnEBHJjw0LdiFathFNOBt8pBhBTX72YJdgismh8PEqg51EjD82OaSgcjWp3Ujn8WmKUMA/h9wEvepY43b3O6DcBmB/U2e6YGYvhvIUxI0ZkGQZsC6zj0r/h8B7jUsr4Fot2cYpnE/GUE+khJifCm0UkIPKlgDStxprXOHCwtgDueDtMow1EtcuAAPJuqEXb7CTaRxVWTZ7iwkYm5tOKmXguVn+MiKy3RxZ41MA9Xd3oOQFHNRQEEN91mP1fNgqMPESzVUmJbBrQPt/+RLQhfMaFJIqtfG7sDeI8sA4itcoj8AKTpUhfhUtuS0q+W30cMY/70i1YBcOgRgzVwfIfECVtO5vm4D4ntdebCcQW3fuj7dCmMtHj04ha9P+7Czeo4LCVLcJ7z4Y0BdrIYV9ZnJYBw3mCK3RLZ99/qGHBesNqrY1SCuU5ZnS1h7T2fdQECmsYYqA2fgvFThKjpB9id/7KNeiWuF3sgD5K7Czv8veWqg4ewCZEqZ7B8BXZZwVYzrFEqsONLu9qiofc6J62A7qT4fM35rfoHNO2kFHgbJF3dYrbcT5y+gZYxN9RKsJuL7uBgATlbqSpELGzKhzXGqCqiJ918Ji+k2hsLbUc6vmSdSV7JCw+6Yr0UfTDd7eFL2VmZyV8/Zy2go4UsXd2Xx+3s22ZbZxUULDdxq1T0/+pIuBHzNLhjxmfI0qQK1NOI\",\"iv\":\"jLsJjzIjzlZvKOnF4lHejQ==\",\"salt\":\"p/6m1lrdIB8G88iwE+cUFO/B7sTDXcJ0ERfzSDIHZZc=\"}";

// Decrypt the payload using the provided password
string encryptedPayload = crypto.Decrypt(password, payload);

// Print the decrypted payload
Console.WriteLine(encryptedPayload);