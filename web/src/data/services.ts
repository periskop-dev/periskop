import { Dispatch } from "redux";
import * as RemoteData from "data/remote-data";
import { registerReducer } from "data/store"
import { ServicesState } from "data/types"

export const FETCH = "periskop/services/FETCH"
export const FETCH_SUCCESS = "periskop/services/FETCH_SUCCESS"
export const FETCH_FAILURE = "periskop/services/FETCH_FAILURE"
export const SET_SERVICE = "periskop/services/SET_SERVICE"

export type ServicesAction =
  | { type: typeof FETCH }
  | { type: typeof FETCH_SUCCESS; services: string[]; defaultService?: string }
  | { type: typeof FETCH_FAILURE; error: any }
  | { type: typeof SET_SERVICE; service: string }

function fetchingServices(): ServicesAction {
  return { type: FETCH }
}

function fetchedServicesSuccessfully(services: string[]): ServicesAction {
  return {
    type: FETCH_SUCCESS,
    services
  }
}

function fetchedServicesFailed(error: any): ServicesAction {
  return { type: FETCH_FAILURE, error}
}

export const fetchServices = (service?: string) => {
  return(
    dispatch: Dispatch<ServicesAction>
  ) => {
    dispatch(fetchingServices())

    return fetch(`${parseHostName()}services/`).then(response => {
      return response
        .json()
        .then(services => dispatch(fetchedServicesSuccessfully(services)))
        .catch(err => dispatch(fetchedServicesFailed(err)))
    })
  }
}

const initialState: ServicesState = {
  services: RemoteData.idle()
}

export default function servicesReducer(state = initialState, action: ServicesAction) {

  switch (action.type) {
    case FETCH:
      return initialState
    case FETCH_SUCCESS:
      const sortedServices = action.services.sort((a: string, b: string) => a.localeCompare(b))
      return {
        ...state,
        services: RemoteData.succeed(sortedServices)
      }
    case FETCH_FAILURE:
      return {
        ...state,
        services: RemoteData.fail(action.error)
      }
    case SET_SERVICE:
      return {
        ...state,
        activeService: action.service
      }
    default:
      return state
  }
}

registerReducer("servicesReducer", servicesReducer)

const parseHostName = () => {
  let windowUrl = new URL(window.location.origin)
  if (windowUrl.hostname === "localhost") {
    windowUrl.port = "7777"
  }
  return windowUrl
}
