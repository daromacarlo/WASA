<template>
	<div class="home">
	  <div class="container">
		<h1 class="title">Benvenuto su Wasachat!</h1>
  
		<div class="button-container">
		  <button class="btn login-btn" @click="goToLogin">Accedi</button>
		</div>
	  </div>
	</div>
  </template>
  
  <script>
  export default {
	name: 'Home',
	methods: {
	  goToLogin() {
		// Naviga alla pagina di login
		this.$router.push('/login');
	  },
	},
  };
  </script>
  
  <style scoped>
  .home {
	background-color: #f7f7f7;
	display: flex;
	justify-content: center;
	align-items: center;
	height: 100vh;
	font-family: 'Arial', sans-serif;
  }
  
  .container {
	text-align: center;
	background-color: white;
	padding: 30px;
	border-radius: 10px;
	box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
  }
  
  .title {
	font-size: 40px;
	color: #333;
	margin-bottom: 10px;
  }
  
  .subtitle {
	font-size: 15px;
	color: #777;
	margin-bottom: 30px;
  }
  
  .button-container {
	display: flex;
	justify-content: center;
	gap: 20px;
  }
  
  .btn {
	padding: 10px 60px;
	border: none;
	border-radius: 5px;
	font-size: 20px;
	cursor: pointer;
	transition: background-color 0.3s ease;
  }
  
  .login-btn {
	background-color: #7dac10;
	color: white;
  }
  
  .login-btn:hover {
	background-color: #7dac10;
  }

  </style>