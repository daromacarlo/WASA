<template>
	<div class="home">
	  <div class="container">
		<h1 class="titolo">WASACHAT</h1>
		<button class="bottone_login" @click="loggati">Accedi</button>
	  </div>
	</div>
  </template>
  
  <script>
  export default {
	name: 'Home',
	methods: {
	  loggati() {
		this.$router.push('/login');
	  },
	},
  };
  </script>
  
  <style scoped>
   .home {
      background-color: #f7e3b8;
      display: flex;
      justify-content: center;
      align-items: center;
      height: 100vh;
      font-family: 'Roboto', sans-serif;
    }
    
    .container {
      text-align: center;
      background-color: rgb(250, 172, 120);
      padding: 40px;
      border-radius: 12px;
      box-shadow: 0 6px 15px rgba(0, 0, 0, 0.1);
    }

	.titolo {
		font-size: 100px;
		color: #333;
		margin-bottom: 20px;
		border: 2px solid #333;
		border-radius: 8px;   
		padding: 10px;
		-webkit-text-stroke: 2px white;
	}

    .bottone_login {
      background-color: rgb(240, 97, 3);
      color: rgb(221, 219, 219);
      padding: 15px 40px;
	  border-radius: 90px;
	  font-size: 30px;
    }
    
  </style>
  