import * as React from "react"
import { withRouter, RouteComponentProps } from "react-router"
import { bindActionCreators, Dispatch, AnyAction } from "redux"
import { connect } from "react-redux"

import { fetchServices } from "data/services"
import { fetchErrors } from "data/errors"
import { StoreState } from "data/types"
import * as RemoteData from "data/remote-data"

import * as moment from "moment"

import { Navbar, NavItem, Nav, NavDropdown } from "react-bootstrap"

interface ConnectedProps {
  services: RemoteData.RemoteData<any, string[]>,
  activeService?: string,
  updatedAt?: number
}

interface DispatchProps {
  fetchServices: (service?: string) => void,
  fetchErrors: (service: string) => void
}

type Props = DispatchProps & ConnectedProps & RouteComponentProps<{service: string}>

class NavbarComponent extends React.Component<Props, {}> {

  constructor(props: Props) {
    super(props)

    this.handleRefreshClick = this.handleRefreshClick.bind(this)
    this.handleServiceSelect = this.handleServiceSelect.bind(this)
  }

  componentDidMount() {
    this.props.fetchServices(this.props.match.params.service)
  }

  handleRefreshClick() {
    this.props.fetchErrors(this.props.activeService)
  }

  handleServiceSelect(service: string) {
    if (this.props.activeService !== service) {
      this.props.history.push(`/${service}`)
    }
  }

  renderServicesInDropdown() {
    if (RemoteData.isSuccess(this.props.services)) {
      return this.props.services.data.map((service, index) => <NavDropdown.Item key={index} onClick={_ => this.handleServiceSelect(service)}>{service}</NavDropdown.Item>)
    }
  }

  renderUpdatedAt() {
    if (this.props.updatedAt !== undefined) {
      return `Updated: ${moment(this.props.updatedAt).calendar()}`
    }
  }

  renderRefreshButton() {
    if (this.props.activeService) {
      return (
        <Nav.Item>
          <button className="btn btn-xs btn-default" onClick={this.handleRefreshClick}>
            Refresh
          </button>
        </Nav.Item>
      )
    }
  }

  render() {
    return (
        <Navbar
          bg="dark"
          collapseOnSelect
          fixed="top"
        >
        <Navbar.Brand>
          <a href="/">Periskop</a>
        </Navbar.Brand>
        <Navbar.Toggle/>
        <Navbar.Collapse>
          <Nav>
            <NavDropdown title={this.props.activeService ? this.props.activeService : "Service"} id="project-nav-dropdown">
            {this.renderServicesInDropdown()}
            </NavDropdown>
          </Nav>
          <Nav className="justify-content-end">
            {this.renderRefreshButton()}
          </Nav>
          <Navbar.Text>{this.renderUpdatedAt()}</Navbar.Text>
        </Navbar.Collapse>
      </Navbar>
    )
  }
}

function matchDispatchToProps(dispatch: Dispatch<AnyAction>): DispatchProps {
  return bindActionCreators({ fetchServices, fetchErrors }, dispatch);
}

function mapStateToProps(state: StoreState): ConnectedProps {
  return {
    services: state.servicesReducer.services,
    updatedAt: state.errorsReducer.updatedAt,
    activeService: state.errorsReducer.activeService
  }
}

export default withRouter(connect<ConnectedProps, DispatchProps, {}>(mapStateToProps, matchDispatchToProps)(NavbarComponent))
