package templates

const TEMPLATE_STRING = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Account Summary</title>
    <style>
      body {
        font-family: "Helvetica Neue", Arial, sans-serif;
        margin: 0;
        padding: 0;
        background-color: #f7f7f7;
        color: #333;
        text-align: center;
      }

      header {
        background-color: #333;
        padding: 10px;
        color: #fff;
      }

      .logo {
        max-width: 100px;
        height: auto;
        margin: 0 auto;
        display: block;
        filter: grayscale(100%);
      }

      .ticket {
        background-color: #fefefe;
        max-width: 400px;
        margin: 20px auto;
        padding: 20px;
        border-radius: 10px;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        position: relative;
        overflow: hidden;
        text-align: left;
      }

      .dashed-line {
        border-bottom: 1px dashed #ccc;
        margin: 10px 0;
      }

      h2 {
        color: #555;
        margin-bottom: 10px;
      }

      h3 {
        color: #777;
        margin-top: 15px;
      }

      p {
        color: #888;
        margin: 5px 0;
      }

      .cut-line {
        background-color: #fff;
        height: 20px;
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        z-index: 1;
        clip-path: polygon(0 100%, 50% 80%, 100% 100%);
        background: linear-gradient(
          90deg,
          rgba(255, 255, 255, 0),
          rgba(255, 255, 255, 1),
          rgba(255, 255, 255, 0)
        );
      }
    </style>
  </head>
  <body>
    <header>
      <img
        src="https://upload.wikimedia.org/wikipedia/commons/e/e3/Stori_logo_vertical.png"
        alt="Logo"
        class="logo"
      />
    </header>
    <div class="ticket">
      <h2>Account Summary</h2>
      <p>Total balance: ${{totalBalance}}</p>
      <div class="dashed-line"></div>
      <h3>Transaction Details</h3>
      {{transactionDetail}}
      <div class="dashed-line"></div>
      <h3>Average Transaction Amounts</h3>
      <p>Average debit amount: -${{avgDebit}}</p>
      <p>Average credit amount: ${{avgCredit}}</p>
      <div class="dashed-line"></div>
    </div>
  </body>
</html>
`
