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



ENV DATABASE_URL=postgresql://postgres:password@postgres:5432/chainlink2?sslmode=disable


ENV FEATURE_OFFCHAIN_REPORTING=true
ENV FEATURE_OFFCHAIN_REPORTING2=true
ENV P2P_NETWORKING_STACK=V1V2
ENV P2PV2_LISTEN_ADDRESSES=0.0.0.0:9901
ENV OCR_KEY_BUNDLE_ID=42290e41b163d00b128f62434344c53ee0fb3d0b4366e5e3193dc8bc4900dbce
ENV P2P_PEER_ID=p2p_12D3KooWRQrtV2chVNoWp3pYtzSaUzBFQJKwy59UDwXsoaaCpE9v
ENV OCR_TRANSMITTER_ADDRESS=0x5E43b2015a2fa5BeECed22626AEC4d688023b2cf



ENV P2P_LISTEN_PORT=9902
#ENV P2P_BOOTSTRAP_PEERS=/dns4/node1/tcp/9900/p2p/12D3KooWJUknXTJfuyFKa7UnJYMSJB811s5ymvBjSQcAjpnSjyBT