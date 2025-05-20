<template>
  <div class="login-view">
    <button class="bottone-uscita" @click="vai_alla_home">Vai alla home</button>
    <div class="login-container">
      <h1 class="titolo">Accedi a WASACHAT</h1>
      <form @submit.prevent="login">
        <input class="input" v-model="nickname" placeholder="Inserisci il tuo nome:"/>
        <button class="bottone">Accedi</button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      nickname: ''
    };
  },

  methods: {
    async vai_alla_home() {
      this.$router.push("/");
      return;
      },

    async login() {
      if (this.nickname === "") {
        alert("Nickname nullo non valido");
        return;
      }

      try {
        const response = await this.$axios.post("/wasachat", { nickname: this.nickname });
        const messaggio = response.data.risposta;
        const codice = parseInt(response.data.codice);

        if (codice === 200) {
          this.$router.push(`/wasachat/${this.nickname}/chats`);
        } else {

          alert(messaggio);
        }
      } catch (e) {
        if (e.response) {
          const messaggio = e.response.data.errore;
          const codice = parseInt(e.response.data.codiceErrore);
          alert(messaggio + ` (codice ${codiceErrore})`);
        } else {
          alert("Errore di rete o server non raggiungibile.");
        }
        console.error(e);
      }
    }
  }
};
</script>

<style scoped>
  .login-view {
    background-color: #f7e3b8;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    font-family: 'Roboto', sans-serif;
  }
  
  .login-container {
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
		padding: 10px;
		-webkit-text-stroke: 2px white;
  }

  .input {
    display: block;
    width: 100%;
    padding: 20px;
    margin: 10px 0;
    border: 1px solid #ccc;
    border-radius: 5px;
  }

  .bottone {
    background-color: rgb(240, 97, 3);
    color: rgb(221, 219, 219);
    padding: 20px 40px;
    margin: 40px ;
	  border-radius: 90px;
	  font-size: 30px;
  }

  .bottone-uscita {
  background-color: rgb(197, 66, 66);
  color: rgb(221, 219, 219);
  padding: 20px 40px;
  margin: 40px;
  border-radius: 90px;
  font-size: 15px;
  
  position: fixed;
  top: 10px;    
  left: 10px;      
}


</style>
