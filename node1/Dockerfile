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

ENV DATABASE_URL=postgresql://postgres:password@postgres:5432/chainlink1?sslmode=disable

ENV FEATURE_OFFCHAIN_REPORTING=true
ENV FEATURE_OFFCHAIN_REPORTING2=true
ENV P2P_NETWORKING_STACK=V1V2
ENV P2PV2_LISTEN_ADDRESSES=0.0.0.0:9900
ENV OCR_KEY_BUNDLE_ID=c2e395733ebddc46ac2c91b0c5fa4a77fdb8923e56bef82e865df89e13ff4542
ENV P2P_PEER_ID=p2p_12D3KooWJUknXTJfuyFKa7UnJYMSJB811s5ymvBjSQcAjpnSjyBT
ENV OCR_TRANSMITTER_ADDRESS=0x6f0e86F85c03E6Da572f59711FfE930bC05fEB69
ENV CHAINLINK_DEV=true

ENV P2P_LISTEN_PORT=7979