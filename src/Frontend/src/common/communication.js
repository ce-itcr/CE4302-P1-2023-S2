import axios from "axios";

export const url = "http://localhost:8080/";

const postRequest = async (
    action, data
) => {

    try {
        const res = await axios({
            method: 'post',
            baseURL: url + action,
            url: url + action,
            headers: { "Content-Type": "application/json" },
            data: data,
        })


        return res.data
    } catch (error) {
        if (error.response) {
            console.log(error.response)
        } else {
            console.log(error.message)
        }
    }
}

const getDashboardData = async (
    params = {}
) => {

    try {
        const res = await axios({
            method: 'get',
            baseURL: url + "about",
            url: url + "about",
            headers: { "Content-Type": "application/json" },
            params: params,
        })


        return res.data
    } catch (error) {
        if (error.response) {
            console.log(error.response)
        } else {
            console.log(error.message)
        }
    }
}

const getMetricsData = async (
    params = {}
) => {

    try {
        const res = await axios({
            method: 'get',
            baseURL: url + "aboutmetrics",
            url: url + "aboutmetrics",
            headers: { "Content-Type": "application/json" },
            params: params,
        })


        return res.data
    } catch (error) {
        if (error.response) {
            console.log(error.response)
        } else {
            console.log(error.message)
        }
    }
}


export { postRequest, getDashboardData, getMetricsData}