import axios from 'axios';

var app = new Vue({
    el: '#app',
    data () {
        return {
            localEvents: null
        }
    },
    mounted () {
        axios
            .get("http://localhost:8081/localevents?postcode=78759&miles=20")
            .then(response => (this.localEvents = response))
    }
})