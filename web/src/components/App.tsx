import * as React from "react"
import { withRouter } from "react-router"
import { connect } from "react-redux"
import { RouteComponentProps } from "react-router"
import { bindActionCreators, Dispatch, AnyAction } from "redux"

import ErrorComponent from "components/Error"
import SideBar from "components/SideBar"

import * as RemoteData from "data/remote-data"
import { StoreState, AggregatedError, SeverityFilter, ErrorsState, Targets } from "data/types"
import { fetchTargets } from "data/targets"
import { fetchErrors, setActiveError, setErrorsSeverityFilter, setErrorsSearchFilter } from "data/errors"

import { Row, Col, Container } from "react-bootstrap"
import { getFilteredErrors } from "util/errors"

interface ConnectedProps {
  errors: RemoteData.RemoteData<any, AggregatedError[]>,
  services: RemoteData.RemoteData<any, string[]>,
  activeError: AggregatedError,
  activeService: string,
  severityFilter: SeverityFilter,
  searchTerm: string,
  targets: RemoteData.RemoteData<any, Targets>,
}

interface DispatchProps {
  fetchErrors: (service: string) => void
  setActiveError: (errorKey: string) => void
  setErrorsSeverityFilter: (severity: SeverityFilter) => void
  setErrorsSearchFilter: (searchTerm: ErrorsState["searchTerm"]) => void
  fetchTargets: () => void
}

type Props = ConnectedProps & DispatchProps & RouteComponentProps<{ service: string, errorKey: string }>


class App extends React.Component<Props, {}> {
  constructor(props, context) {
    super(props, context)
    this.handlerErrorSelect = this.handlerErrorSelect.bind(this)
  }

  componentDidMount() {
    if (RemoteData.isSuccess(this.props.services)) {
      if (this.props.match.params.service !== "targets") {
        this.props.fetchErrors(this.props.match.params.service)
      }
      this.props.fetchTargets()
    }
  }

  componentDidUpdate(prevProps: Props) {
    const { activeService, services, errors, match, fetchErrors, setActiveError, fetchTargets } = this.props

    if (RemoteData.isSuccess(services)) {
      if (services.data.includes(match.params.service)
        && (activeService !== match.params.service)
        && !RemoteData.isLoading(errors)) {
          if (match.params.service !== "targets") {
            fetchErrors(match.params.service)
          }
          fetchTargets()
      }
    }

    if (RemoteData.isSuccess(errors)) {
      const decodedErrorKey = decodeURIComponent(match.params.errorKey)
      const activeError = errors.data.find(e => e.aggregation_key === decodedErrorKey)

      if (!activeError || RemoteData.isLoading(errors)) return

      const hasNewActiveError = this.props.activeError?.aggregation_key !== activeError.aggregation_key
      if (hasNewActiveError) {
        setActiveError(activeError.aggregation_key)
      }
    }
  }

  handlerErrorSelect(errorKey: string) {
    this.props.history.push(`/${this.props.match.params.service}/errors/${encodeURIComponent(errorKey)}`)
  }

  renderSideBar() {
    switch (this.props.errors.status) {
      case RemoteData.SUCCESS:
        return (
          <SideBar
            errors={this.props.errors.data}
            handleErrorSelect={this.handlerErrorSelect}
          />
        )
      case RemoteData.LOADING:
        return <div>fetching errors...</div>
    }
  }

  renderError() {
    switch (this.props.errors.status) {
      case RemoteData.SUCCESS:
        if ((this.props.errors.data.length !== 0 && this.props.activeError !== undefined)) {
          return <ErrorComponent />
        }
    }
  }

  render() {
    return (
      <div className="app-component">
        <Container fluid className="app-component-grid">
          <Row className="app-component-row">
            <Col xs={3} id="left-column">
              {this.renderSideBar()}
            </Col>
            <Col xs={9} id="right-column">
              {this.renderError()}
            </Col>
          </Row>
        </Container>
      </div>
    )
  }
}


const mapStateToProps = ({ errorsReducer, servicesReducer, targetsReducer }: StoreState) => {
  const { errors, activeError, severityFilter, searchTerm, activeService, activeSortFilter } = errorsReducer
  const { services } = servicesReducer
  const { targets } = targetsReducer

  const defaultConnectedProps = { activeError, services, severityFilter, activeService, searchTerm, targets }

  if (RemoteData.isSuccess(errors)) {
    const filteredErrors = getFilteredErrors(errors.data, searchTerm, severityFilter, activeSortFilter)

    return {
      ...defaultConnectedProps,
      errors: RemoteData.succeed(filteredErrors),
    }
  }

  return { ...defaultConnectedProps, errors }
}


const matchDispatchToProps = (dispatch: Dispatch<AnyAction>): DispatchProps => {
  return bindActionCreators({ fetchErrors, setActiveError, setErrorsSeverityFilter, setErrorsSearchFilter, fetchTargets }, dispatch);
}

export default withRouter(connect<ConnectedProps, {}, RouteComponentProps<{ service: string }>>(mapStateToProps, matchDispatchToProps)(App))
