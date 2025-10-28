import request from './request'

export function getCloudVars(params) {
  return request({
    url: '/admin/cloud-vars',
    method: 'get',
    params
  })
}

export function setCloudVar(data) {
  return request({
    url: '/admin/cloud-vars',
    method: 'post',
    data
  })
}

export function deleteCloudVar(id) {
  return request({
    url: `/admin/cloud-vars/${id}`,
    method: 'delete'
  })
}

export function batchSetCloudVars(data) {
  return request({
    url: '/admin/cloud-vars/batch',
    method: 'post',
    data
  })
}

export function batchDeleteCloudVars(data) {
  return request({
    url: '/admin/cloud-vars/batch',
    method: 'delete',
    data
  })
}

