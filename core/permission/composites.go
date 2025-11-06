package permission

import (
	"grf/core/exceptions"

	"github.com/gofiber/fiber/v2"
)

type And struct {
	Perms []IPermission
}

func NewAnd(perms ...IPermission) IPermission {
	return &And{Perms: perms}
}

func (a *And) Check(c *fiber.Ctx) error {
	for _, perm := range a.Perms {
		if err := perm.Check(c); err != nil {
			return err
		}
	}
	return nil
}

type Or struct {
	Perms []IPermission
}

func NewOr(perms ...IPermission) IPermission {
	return &Or{Perms: perms}
}

func (o *Or) Check(c *fiber.Ctx) error {
	for _, perm := range o.Perms {
		if err := perm.Check(c); err == nil {
			return nil
		}
	}

	return exceptions.NewForbidden("permission_denied", nil)
}

type Not struct {
	Perm IPermission
}

func NewNot(perm IPermission) IPermission {
	return &Not{Perm: perm}
}

func (n *Not) Check(c *fiber.Ctx) error {
	err := n.Perm.Check(c)
	if err == nil {
		return exceptions.NewForbidden("permission_denied", nil)
	}
	return nil
}
