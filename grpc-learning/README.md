# Mini Stock Exchange with gRPC

This is a simplified stock trading platform.
Users can deposit or withdraw funds in dollars, use them to buy or sell shares of an imaginary stock called **AwesomeStock**, and request statements of their trades and account activity.
The platform also provides real-time stock price updates and supports different gRPC communication patterns.

## Use-Cases and gRPC Routines
### 1. DepositFunds, WithdrawFunds, GetBalance
- Type: Request–Response (Unary RPC)
- Description: Users can deposit or withdraw dollars into their trading account and check their current balance.

### 2. TradeStatement
- Type: Client-Side Streaming
- Description: The client streams all trade events (buy/sell actions) of a user to the statement service. The statement service then processes them and returns a consolidated statement summarizing the user’s trades.

### 3. GetTradeHistory
- Type: Server-Side Streaming
- Description: The client requests a user’s trading history. The service responds by streaming all trades back, one by one, in chronological order.

### 4. LiveTrading (Market Data + Orders)
- Type: Bidirectional Streaming
- Description: The server streams real-time price updates for **AwesomeStock**. The client can send trading commands (Buy, Sell) in response to market updates.
The server processes the trades immediately and confirms execution results over the same stream.
