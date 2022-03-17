FROM smartcontract/chainlink:1.2.0

# Create the chainlink node root path
RUN mkdir /chainlink

# Arg for api user email, with default value
ARG API_USER_EMAIL="user@example.com"
ENV API_USER_EMAIL=$API_USER_EMAIL

# Arg for api user password, with default value
ARG API_USER_PASSWORD="PA@SSword1234!567"
ENV API_USER_PASSWORD=$API_USER_PASSWORD

# Arg for node wallet password, with default value
ARG WALLET_PASSWORD="PA@SSword1234!567"
ENV WALLET_PASSWORD=$WALLET_PASSWORD

# Create chainlink node required values to initialize with
RUN echo $API_USER_EMAIL > /chainlink/.api
RUN echo $API_USER_PASSWORD >> /chainlink/.api
RUN echo $WALLET_PASSWORD > /chainlink/.password

ENV FEATURE_OFFCHAIN_REPORTING=true
ENV FEATURE_OFFCHAIN_REPORTING2=true

ENV DATABASE_URL=postgresql://postgres:password@postgres:5432/chainlink3?sslmode=disable

# ENV OCR_KEY_BUNDLE_ID=298de9c4b923653c881d95c125ecef49a2c6a075c9ac9356cce1ccf81218aecc
# ENV P2P_PEER_ID=p2p_12D3KooWP5bFi2CQvHHKSUrmtqF7YumDFE9UNR4x6fs9DpanKQLv
# ENV OCR_TRANSMITTER_ADDRESS=0x846a5B0838e20bB892B943Ac25F60A01c154b451


ENV P2P_LISTEN_PORT=7979
#ENV P2P_BOOTSTRAP_PEERS=/dns4/node1/tcp/7777/p2p/12D3KooWRiWTNQTnHMgaezYXjEeYFiRJxJjj2KneBsW6o8qXkwvB