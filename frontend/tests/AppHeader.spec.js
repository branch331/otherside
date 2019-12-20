import { shallowMount } from '@vue/test-utils'
import AppHeader from '../src/components/AppHeader.vue'
import Vue from 'vue'
import { isIterable } from 'core-js';
import { shallow } from 'vue-test-utils';

describe('AppHeader', () => {
    
    it('calls getCurrentUser if token and apiAddress are provided', () => {
        const getCurrentUser = jest.fn(x => "Real User")
        const wrapper = shallowMount(AppHeader, {
            propsData: {
                token: "totallyRealToken3453435",
                apiAddress: "superRealEndpoint"
            },
            methods: {
                getCurrentUser
            }
        })

        expect(getCurrentUser).toBeCalled()
    })

    it('does not call getCurrentUser if token is missing', () => {
        const getCurrentUser = jest.fn(x => "Real User")
        const wrapper = shallowMount(AppHeader, {
            propsData: {
                token: null,
                apiAddress: "superRealEndpoint"
            },
            methods: {
                getCurrentUser
            }
        })

        expect(getCurrentUser).not.toBeCalled()
    })

    it('does not call getCurrentUser if apiAddress is missing', () => {
        const getCurrentUser = jest.fn(x => "Real User")
        const wrapper = shallowMount(AppHeader, {
            propsData: {
                token: "waterbottle59",
                apiAddress: null
            },
            methods: {
                getCurrentUser
            }
        })

        expect(getCurrentUser).not.toBeCalled()
    })

    it('does not call getCurrentUser if currentUser is already populated', async () => {
        const getCurrentUser = jest.fn(x => "Real User")
        const wrapper = shallow(AppHeader, {
            propsData: {
                token: "superSpecialToken",
                apiAddress: "myapiaddress"
            },
            data: function() {
                return {
                    currentUser: "Johnny boy"
                }
            },
            methods: {
                getCurrentUser
            }
        })

        expect(getCurrentUser).not.toBeCalled()
    })
})