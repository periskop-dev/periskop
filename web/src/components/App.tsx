import * as React from "react"
import { withRouter } from "react-router"
import { connect } from "react-redux"
import { RouteComponentProps } from "react-router"
import { bindActionCreators, Dispatch, AnyAction } from "redux"

import ErrorComponent from "components/Error"
import SideBar from "components/SideBar"

import * as RemoteData from "data/remote-data"
import { StoreState, AggregatedError } from "data/types"
import { fetchErrors, setActiveError } from "data/errors"

import { Row, Col, Container } from "react-bootstrap"
import { filterErrorsBySubstringMatch } from "util/errors"

interface ConnectedProps {
  errors: RemoteData.RemoteData<any, AggregatedError[]>,
  services: RemoteData.RemoteData<any, string[]>,
  activeError: AggregatedError,
  activeService: string,
}

interface DispatchProps {
  fetchErrors: (service: string) => void
  setActiveError: (errorKey: string) => void
}

type Props = ConnectedProps & DispatchProps & RouteComponentProps<{ service: string, errorKey: string }>

export interface State {
  errors: AggregatedError[]
  searchTerm: string,
}


class App extends React.Component<Props, State> {
  state = {
    errors: [],
    searchTerm: "",
  }

  constructor(props, context) {
    super(props, context)
    this.handlerErrorSelect = this.handlerErrorSelect.bind(this)
  }

  componentDidMount() {
    if (RemoteData.isSuccess(this.props.services)) {
      this.props.fetchErrors(this.props.match.params.service)
    }
  }

  componentDidUpdate(prevProps: Props) {
    const { activeService, services, errors, match, fetchErrors, setActiveError } = this.props
    const { searchTerm } = this.state

    if (RemoteData.isSuccess(services)) {
      if (
        services.data.includes(match.params.service) &&
        (activeService !== match.params.service) &&
        !RemoteData.isLoading(errors)) {
        fetchErrors(match.params.service)
      }
    }

    if (RemoteData.isSuccess(errors)) {
      let decodedErrorKey = decodeURIComponent(match.params.errorKey)
      let activeError = errors.data.find(e => e.aggregation_key === decodedErrorKey)
      if (
        activeError !== undefined &&
        (activeError !== decodedErrorKey) &&
        !RemoteData.isLoading(errors)) {
        setActiveError(activeError.aggregation_key)
      }

      const hasNewErrors = errors !== prevProps.errors && errors
      if (hasNewErrors) {
        this.handleFilterByAggregatedKey(searchTerm)
      }
    }
  }

  handlerErrorSelect(errorKey: string) {
    this.props.history.push(`/${this.props.match.params.service}/errors/${encodeURIComponent(errorKey)}`)
  }

  handleFilterByAggregatedKey = (searchTerm: string) => {
    const { errors } = this.props

    switch (errors.status) {
      case RemoteData.SUCCESS:
        return this.setState({
          errors: filterErrorsBySubstringMatch(errors.data, searchTerm),
          searchTerm,
        })

      case RemoteData.LOADING:
        return <div>fetching errors...</div>
    }
  }

  renderSideBar() {
    switch (this.props.errors.status) {
      case RemoteData.SUCCESS:
        if (this.props.errors.data.length === 0) {
          return <div>no errors returned by api</div>
        } else {
          return (
            <SideBar
              errors={this.state.errors}
              handleErrorSelect={this.handlerErrorSelect}
              onSearchByAggredgatedKey={this.handleFilterByAggregatedKey}
              searchKey={this.state.searchTerm}
            />
          )
        }
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

const mapStateToProps = (state: StoreState) => {
  return {
    errors: state.errorsReducer.errors,
    activeError: state.errorsReducer.activeError,
    services: state.servicesReducer.services,
    activeService: state.errorsReducer.activeService
  }
}

const matchDispatchToProps = (dispatch: Dispatch<AnyAction>): DispatchProps => {
  return bindActionCreators({ fetchErrors, setActiveError }, dispatch);
}

export default withRouter(connect<ConnectedProps, {}, RouteComponentProps<{ service: string }>>(mapStateToProps, matchDispatchToProps)(App))
