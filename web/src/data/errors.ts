import { Dispatch } from "redux";
import * as RemoteData from "data/remote-data";
import { registerReducer } from "data/store"
import { AggregatedError, ErrorsState } from "data/types"

export const FETCH = "periskop/errors/FETCH"
export const FETCH_SUCCESS = "periskop/errors/FETCH_SUCCESS"
export const FETCH_FAILURE = "periskop/errors/FETCH_FAILURE"
export const SET_ACTIVE_ERROR = "periskop/errors/SET_ACTIVE_ERROR"

export type ErrorsAction =
  | { type: typeof FETCH; service: string }
  | { type: typeof FETCH_SUCCESS; errors: AggregatedError[] }
  | { type: typeof FETCH_FAILURE; error: any }
  | { type: typeof SET_ACTIVE_ERROR; errorKey: string }

export const fetchErrors = (service: string) => {
  return(
    dispatch: Dispatch<ErrorsAction>
  ) => {
    dispatch(fetchingErrors(service))

    return fetch(`${parseHostName()}services/${service}/errors/`).then(response => {
      return response
        .json()
        .then(errors => dispatch(fetchedErrorsSuccessfully(errors, service)))
        .catch(err => dispatch(fetchedErrorsFailed(err)))
    })
  }
}

export const setActiveError = (errorKey: string) => {
  return { type: SET_ACTIVE_ERROR, errorKey }
}

export function fetchingErrors(service: string): ErrorsAction {
  return { type: FETCH, service }
}

export function fetchedErrorsSuccessfully(errors: AggregatedError[], service: string): ErrorsAction {
  return { type: FETCH_SUCCESS, errors}
}

export function fetchedErrorsFailed(error: any): ErrorsAction {
  return { type: FETCH_FAILURE, error }
}

const initialState: ErrorsState = {
  errors: RemoteData.idle(),
  activeError: undefined,
  updatedAt: undefined
}

function errorsReducer(state = initialState, action: ErrorsAction) {
  console.log(action)
  switch (action.type) {
    case FETCH:
      return {
        ...state,
        errors: RemoteData.load(),
        activeService: action.service
      }
    case FETCH_SUCCESS:
      const sortedErrors = action.errors.sort((a: AggregatedError, b: AggregatedError) => b.latest_errors[0].timestamp - a.latest_errors[0].timestamp)
      return {
        ...state,
        errors: RemoteData.succeed(sortedErrors),
        activeError: undefined,
        updatedAt: (new Date()).getTime()
      }
    case FETCH_FAILURE:
      return {
        errors: RemoteData.fail(action.error)
      }
    case SET_ACTIVE_ERROR:
      switch (state.errors.status) {
        case RemoteData.SUCCESS:
          return {
            ...state,
            activeError: state.errors.data.find(e => e.aggregation_key == action.errorKey)
          }
        default: {
          return state
        }
      }
    default:
      return state
  }
}

registerReducer("errorsReducer", errorsReducer)

const parseHostName = () => {
  let windowUrl = new URL(window.location.origin)
  if (windowUrl.hostname === "localhost") {
    windowUrl.port = "7777"
  }
  return windowUrl
}
