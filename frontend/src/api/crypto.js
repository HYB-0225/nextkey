import axios from 'axios'
import request from './request'

export function getEncryptionSchemes() {
  return axios.get('/api/crypto/schemes').then(res => res.data)
}

export function updateProjectEncryption(id, data) {
  return request({
    url: `/admin/projects/${id}/encryption`,
    method: 'post',
    data
  })
}

