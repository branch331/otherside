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
        <b-carousel-slide class="performerSlide" v-for="(performer, i) in performers" v-bind:key="`${i}-${performer}-${carouselId}`" img-src="https://picsum.photos/1024/1024/?image=54">
            <p>{{performer}}</p>
        </b-carousel-slide>
    </b-carousel>
</template>

<script>
export default {
    name: 'PerformerCarousel',
    props: {
        carouselId: String,
        performers: Array
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