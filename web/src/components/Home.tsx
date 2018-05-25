import * as React from "react"
import { RouteComponentProps, withRouter } from "react-router";
import * as RemoteData from "data/remote-data"
import { bindActionCreators, Dispatch, AnyAction } from "redux";
import { StoreState } from "data/types";
import { connect } from "react-redux";
import { fetchServices } from "data/services";
import { Grid, Row, Col, ListGroup, ListGroupItem, Well } from "react-bootstrap";

interface ConnectedProps {
  services: RemoteData.RemoteData<any, string[]>
}

interface DispatchProps {
  fetchServices: (service?: string) => void
}

type Props = DispatchProps & ConnectedProps & RouteComponentProps<{service: string}>

class Home extends React.Component<Props, {}> {

  constructor(props: Props) {
    super(props)

    this.handleServiceSelect = this.handleServiceSelect.bind(this)
  }

  componentDidMount() {
    this.props.fetchServices(this.props.match.params.service)
  }

  handleServiceSelect(service: string) {
    this.props.history.push(`/${service}`)
  }

  renderService(service: string, index: number) {
    return (
      <ListGroupItem key={index} onClick={_ => this.handleServiceSelect(service)}>
        {service}
      </ListGroupItem>
    )
  }

  renderServices() {
    if (RemoteData.isSuccess(this.props.services)) {
      return (
        <ListGroup>
          { this.props.services.data.map((service, index) => this.renderService(service, index)) }
        </ListGroup>
      )
    }
  }

  render() {
    return (
      <div>
          <Grid fluid>
            <Row className="show-grid">
              <Col xs={8} xsOffset={2}>
              <Well>
                <p>
                  <strong>Select a service</strong>
                </p>
                {this.renderServices()}
              </Well>
              </Col>
            </Row>
          </Grid>
        </div>
    )
  }
}

function matchDispatchToProps(dispatch: Dispatch<AnyAction>): DispatchProps {
  return bindActionCreators({ fetchServices }, dispatch);
}

function mapStateToProps(state: StoreState): ConnectedProps {
  return {
    services: state.servicesReducer.services
  }
}

export default withRouter(connect<ConnectedProps, DispatchProps, {}>(mapStateToProps, matchDispatchToProps)(Home))
