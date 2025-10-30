import request from './request'

export function getEncryptionSchemes() {
  return request({
    url: '/api/crypto/schemes',
    method: 'get'
  })
}

export function updateProjectEncryption(id, data) {
  return request({
    url: `/admin/projects/${id}/encryption`,
    method: 'post',
    data
  })
}

