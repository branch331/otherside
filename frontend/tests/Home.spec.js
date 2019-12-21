import { createLocalVue, shallowMount, shallow } from '@vue/test-utils'
import Vuex from 'vuex'
import Vue from 'vue'
import Home from '../src/components/Home.vue'

const localVue = createLocalVue()
localVue.use(Vuex)

const $route = {
    fullPath: 'myFullPath',
    query: 'myquery'
}

describe('Home', () => {
    let wrapper

    it('calls getlocalevents if selectedCities and selectedGenres are both nonzero', () => {
        const getLocalEvents = jest.fn(x => "Real User")

        wrapper = shallowMount(Home, {
            mocks: {
                $route
            },
            computed: {
                selectedCities() {
                    return ['Amsterdam']
                },
                selectedGenres() {
                    return ['Sad Rock', 'Elegant Jazz']
                }
            },
            methods: {
                getLocalEvents
            }
        })

        expect(getLocalEvents).toBeCalled()
    })

    it('displays login button if not logged in, and build playlist button if logged in', async () => {
        wrapper = shallowMount(Home, {
            mocks: {
                $route
            }
        })

        wrapper.vm.isStateStringCorrect = false
        expect(wrapper.find('button').text()).toEqual('Log In')

        wrapper.vm.isStateStringCorrect = true
        await Vue.nextTick()
        expect(wrapper.find('button').text()).toEqual('Build Playlist')
    })
})

