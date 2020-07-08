package ros

// IMPORT REQUIRED PACKAGES.

import (
	"strings"

	"github.com/edwinhayes/rosgo/libgengo"
	"github.com/pkg/errors"
)

// DEFINE PUBLIC STRUCTURES.

type DynamicServiceType struct {
	name    string
	md5sum  string
	text    string
	reqType MessageType
	resType MessageType
	//spec *libgengo.SrvSpec - This may be less convenient?
}

type DynamicService struct {
	dynamicType *DynamicServiceType
	Request     Message
	Response    Message
}

func NewDynamicServiceType(typeName string) (*DynamicServiceType, error) {
	return newDynamicServiceTypeNested(typeName, "")
}

func newDynamicServiceTypeNested(typeName string, packageName string) (*DynamicServiceType, error) {
	// Create an empty message type.
	m := new(DynamicServiceType)

	// Create a message context if for some reason it does not exist yet, as it also contains service definitions
	if context == nil {
		// Create context for our ROS install.
		c, err := libgengo.NewMsgContext(strings.Split(GetRuntimePackagePath(), ":"))
		if err != nil {
			return nil, err
		}
		context = c
	}
	// We need to try to look up the full name, in case we've just been given a short name.
	fullname := typeName

	_, ok := context.GetMsgs()[fullname]
	if !ok {
		// Seems like the package_name we were give wasn't the full name.

		// Messages in the same package are allowed to use relative names, so try prefixing the package.
		if packageName != "" {
			fullname = packageName + "/" + fullname
		}
	}

	// Load context for the target message.
	spec, err := context.LoadSrv(fullname)
	if err != nil {
		return nil, err
	}

	// Now we know all about the service!
	m.name = spec.ShortName
	m.md5sum = spec.MD5Sum
	m.text = spec.Text
	m.reqType, err = NewDynamicMessageType(spec.Request.FullName)
	if err != nil {
		return nil, errors.Wrap(err, "error generating request type")
	}
	m.resType, err = NewDynamicMessageType(spec.Response.FullName)
	if err != nil {
		return nil, errors.Wrap(err, "error generating request type")
	}

	// We've successfully made a new service type matching the requested ROS type.
	return m, nil

}

func (t *DynamicServiceType) MD5Sum() string {
	return t.md5sum
}

func (t *DynamicServiceType) Name() string {
	return t.name
}

func (t *DynamicServiceType) NewService() Service {
	// Don't instantiate services for incomplete types.
	if t == nil {
		return nil
	}
	// But otherwise, make a new one.
	d := new(DynamicService)
	d.dynamicType = t
	d.Request = t.RequestType().NewMessage()
	d.Response = t.ResponseType().NewMessage()
	return d
}

func (t *DynamicServiceType) RequestType() MessageType {
	return t.reqType
}

func (t *DynamicServiceType) ResponseType() MessageType {
	return t.resType
}

func (s *DynamicService) ReqMessage() Message {
	return s.Request
}

func (s *DynamicService) ResMessage() Message {
	return s.Response
}

func (s *DynamicService) Type() ServiceType {
	return s.dynamicType
}
