import axios from 'axios'
import { useContext } from 'react'
import { AuthContext } from '../context/AuthContext'

export const useApi = () => {
  const { token } = useContext(AuthContext)

  const client = axios.create({
    baseURL: process.env.REACT_APP_API_BASE || 'http://localhost:8080',
    headers: {
      Accept: 'application/json',
    },
  })

  client.interceptors.request.use((config) => {
    if (token) config.headers['Authorization'] = `Bearer ${token}`
    return config
  })

  return client
}
