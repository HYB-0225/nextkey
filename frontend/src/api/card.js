import request from './request'

export function getCards(params) {
  return request({
    url: '/admin/cards',
    method: 'get',
    params
  })
}

export function createCards(data) {
  return request({
    url: '/admin/cards',
    method: 'post',
    data
  })
}

export function getCard(id) {
  return request({
    url: `/admin/cards/${id}`,
    method: 'get'
  })
}

export function updateCard(id, data) {
  return request({
    url: `/admin/cards/${id}`,
    method: 'put',
    data
  })
}

export function deleteCard(id) {
  return request({
    url: `/admin/cards/${id}`,
    method: 'delete'
  })
}

export function batchUpdateCards(data) {
  return request({
    url: '/admin/cards/batch',
    method: 'put',
    data
  })
}

export function batchDeleteCards(data) {
  return request({
    url: '/admin/cards/batch',
    method: 'delete',
    data
  })
}

export function freezeCard(id) {
  return request({
    url: `/admin/cards/${id}/freeze`,
    method: 'put'
  })
}

export function unfreezeCard(id) {
  return request({
    url: `/admin/cards/${id}/unfreeze`,
    method: 'put'
  })
}

export function batchFreezeCards(data) {
  return request({
    url: '/admin/cards/batch/freeze',
    method: 'put',
    data
  })
}

export function batchUnfreezeCards(data) {
  return request({
    url: '/admin/cards/batch/unfreeze',
    method: 'put',
    data
  })
}

