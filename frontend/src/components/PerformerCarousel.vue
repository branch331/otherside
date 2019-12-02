<template>
    <b-carousel
        :id="carouselId"
        class="performerCarousel"
        ref="carousel"
        :interval="1000"
        :fade=true
        :no-hover-pause=true 
        :no-touch=true 
        img-width="1024"
        img-height="480"
        style="text-shadow: 1px 1px 2px #333;">
        <b-carousel-slide class="performerSlide" v-for="(artist, i) in artistImages" v-bind:key="`${i}-${artist}-${carouselId}`" :img-src=artist.ImageURL>
            <p>{{artist.Name}}</p>
        </b-carousel-slide>
    </b-carousel>
</template>

<script>
export default {
    name: 'PerformerCarousel',
    props: {
        carouselId: String,
        artistImages: Array
    },
    methods: {
        pauseWrapper: function() {
            this.$refs.carousel.pause();
        },
        setSlideWrapper: function(index) {
            this.$refs.carousel.setSlide(index);
        },
        startWrapper: function() {
            this.$refs.carousel.start();
        }
    }
}
</script>