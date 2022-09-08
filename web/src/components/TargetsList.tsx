import * as React from "react"
import { withRouter, RouteComponentProps } from "react-router"
import * as RemoteData from "data/remote-data"
import { bindActionCreators, Dispatch, AnyAction } from "redux";
import { StoreState } from "data/types";
import { Targets, Target } from "data/types";
import { connect } from "react-redux";
import { fetchTargets } from "data/targets";
import { Row, Col, ListGroup, Accordion, Card, Container } from "react-bootstrap";

interface ConnectedProps {
  targets: RemoteData.RemoteData<any, Targets>
}

interface DispatchProps {
  fetchTargets: () => void
}

type Props = DispatchProps & ConnectedProps & RouteComponentProps<{service: string}>

class TargetsList extends React.Component<Props, {}> {

  componentDidMount() {
    this.props.fetchTargets()
  }

  renderTarget(targetName: string, targets: Target[], index: number) {
    const hosts = targets.map(
      (target, index) =>
        <ListGroup.Item key={targetName+String(index)}>
          <a href={`http://${target.endpoint}`} target="_blank" rel="noreferrer">
            {target.endpoint}
          </a>
        </ListGroup.Item>
    )
    const accordionBody = targets.length > 0 ? hosts : "No endpoints in target"
    return (
      <Card key={targetName}>
        <Card.Header style={{cursor:"pointer"}}>
          <Accordion.Toggle as={ListGroup} eventKey={String(index)}>
            {targetName}
          </Accordion.Toggle>
        </Card.Header>
        <Accordion.Collapse eventKey={String(index)}>
          <Card.Body>{accordionBody}</Card.Body>
        </Accordion.Collapse>
      </Card>
    )
  }

  renderTargets() {
    if (RemoteData.isSuccess(this.props.targets)) {
      const targets = this.props.targets.data
      return (
        <Accordion>
          { Object.keys(targets).map((targetName, index) => this.renderTarget(targetName, targets[targetName], index)) }
        </Accordion>
      )
    }
  }

  render() {
    return (
      <div>
          <Container fluid>
            <Row className="show-grid">
              <Col md={{ span: 8, offset: 2 }}>
                <Card bg="light">
                  <Card.Header as="h5">Scraping targets</Card.Header>
                  <Card.Body>
                    {this.renderTargets()}
                  </Card.Body>
                </Card>
              </Col>
            </Row>
          </Container>
        </div>
    )
  }
}

function matchDispatchToProps(dispatch: Dispatch<AnyAction>): DispatchProps {
  return bindActionCreators({ fetchTargets }, dispatch);
}

function mapStateToProps(state: StoreState): ConnectedProps {
  return {
    targets: state.targetsReducer.targets
  }
}

export default withRouter(connect(mapStateToProps, matchDispatchToProps)(TargetsList))
