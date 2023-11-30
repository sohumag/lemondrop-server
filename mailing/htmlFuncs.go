package mailing

import "fmt"

func GetMailingListEmail() (string, error) {
	return fmt.Sprintf(`
	<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Document</title>
		<link
			href="https://fonts.googleapis.com/css?family=Lobster"
			rel="stylesheet"
		/>
		<link
			href="https://fonts.googleapis.com/css?family=Inter"
			rel="stylesheet"
		/>
	</head>
	<body>
		<div>
			<h2 class="logo" style="color: #5b40f6; margin-bottom: 2em">
				lemondrop
			</h2>

			<!-- <h3>Hey!</h3> -->
			<h1>Glad You Could Join Us On This Trip..</h1>

			<p>
				Congrats! You made it to Lemondrop. We will be launching very
				soon and you will be the first to know all the details about how
				to start your sports entertainment journey. There's nothing for
				you to do right now. We'll let you know when its time to start.
			</p>

			<!-- <div class="btn-holder">
				<button>Click Here to Verify </button>
			</div> -->
		</div>

		<style>
			* {
				font-family: 'Inter';
			}

			.logo {
				font-family: 'Lobster';
				font-weight: bold;
			}

			.btn-holder {
				width: 100%;
				display: flex;
				justify-content: center;
				align-items: center;
			}

			button {
				border-radius: 20px;
				padding: 1em 2em;
				font-weight: bold;
				background-color: #5b40f6;
				color: white;
				border: none;
				margin: 2em;
			}
		</style>
	</body>
</html>
`), nil
}

func GetEmailVerificationEmail(email string) (string, error) {
	return fmt.Sprintf(`
	<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Document</title>
		<link
			href="https://fonts.googleapis.com/css?family=Lobster"
			rel="stylesheet"
		/>
		<link
			href="https://fonts.googleapis.com/css?family=Inter"
			rel="stylesheet"
		/>
	</head>
	<body>
		<div>
			<h2 class="logo" style="color: #5b40f6; margin-bottom: 2em">
				lemondrop
			</h2>

			<h3>Hey!</h3>
			<h1>Verify Your Email Address</h1>

			<p>
				As an extra security measure, please verify this is the correct
				email address for your Lemondrop account.
			</p>

			<div class="btn-holder">
				<button>Click Here to Verify</button>
			</div>
		</div>

		<style>
			* {
				font-family: 'Inter';
			}

			.logo {
				font-family: 'Lobster';
				font-weight: bold;
			}

			.btn-holder {
				width: 100%;
				display: flex;
				justify-content: center;
				align-items: center;
			}

			button {
				border-radius: 20px;
				padding: 1em 2em;
				font-weight: bold;
				background-color: #5b40f6;
				color: white;
				border: none;
				margin: 2em;
			}
		</style>
	</body>
</html>


	
	`), nil
}
