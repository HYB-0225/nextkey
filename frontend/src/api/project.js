import request from './request'

export function getProjects(params) {
  return request({
    url: '/admin/projects',
    method: 'get',
    params
  })
}

export function createProject(data) {
  return request({
    url: '/admin/projects',
    method: 'post',
    data
  })
}

export function updateProject(id, data) {
  return request({
    url: `/admin/projects/${id}`,
    method: 'put',
    data
  })
}

export function deleteProject(id) {
  return request({
    url: `/admin/projects/${id}`,
    method: 'delete'
  })
}

export function getProjectByUUID(uuid) {
  return request({
    url: `/admin/projects/${uuid}`,
    method: 'get'
  })
}

export function batchCreateProjects(data) {
  return request({
    url: '/admin/projects/batch',
    method: 'post',
    data
  })
}

export function batchDeleteProjects(data) {
  return request({
    url: '/admin/projects/batch',
    method: 'delete',
    data
  })
}

