package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	mongoModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func UsersRepositoryToUsersMapper(usersRepository UsersRepository) userModel.Users {
	users := make([]userModel.User, len(usersRepository.Users))
	for index, userRepository := range usersRepository.Users {
		users[index] = UserRepositoryToUserMapper(userRepository)
	}

	return userModel.NewUsers(
		users,
		usersRepository.PaginationResponse,
	)
}

func UserRepositoryToUserMapper(userRepository UserRepository) userModel.User {
	return userModel.NewUser(
		userRepository.ID.Hex(),
		userRepository.CreatedAt,
		userRepository.UpdatedAt,
		userRepository.Name,
		userRepository.Email,
		userRepository.Password,
		userRepository.Role,
		userRepository.Verified,
	)
}

func UserCreateToUserCreateRepositoryMapper(userCreate userModel.UserCreate) UserCreateRepository {
	return NewUserCreateRepository(
		userCreate.Name,
		userCreate.Email,
		userCreate.Password,
		userCreate.Role,
		userCreate.Verified,
		userCreate.VerificationCode,
		userCreate.CreatedAt,
		userCreate.UpdatedAt,
	)
}

func UserUpdateToUserUpdateRepositoryMapper(logger applicationModel.Logger, location string, userUpdate userModel.UserUpdate) commonModel.Result[UserUpdateRepository] {
	userObjectID := mongoModel.HexToObjectIDMapper(logger, location+".UserUpdateToUserUpdateRepositoryMapper", userUpdate.ID)
	if validator.IsError(userObjectID.Error) {
		return commonModel.NewResultOnFailure[UserUpdateRepository](userObjectID.Error)
	}

	return commonModel.NewResultOnSuccess(NewUserUpdateRepository(
		userObjectID.Data,
		userUpdate.Name,
		userUpdate.UpdatedAt,
	))
}

func UserResetExpiryRepositoryToUserResetExpiryMapper(userResetExpiry UserResetExpiryRepository) userModel.UserResetExpiry {
	return userModel.NewUserResetExpiry(
		userResetExpiry.ResetExpiry,
	)
}

func UserForgottenPasswordToUserForgottenPasswordRepositoryMapper(userForgottenPassword userModel.UserForgottenPassword) UserForgottenPasswordRepository {
	return NewUserForgottenPasswordRepository(
		userForgottenPassword.ResetToken,
		userForgottenPassword.ResetExpiry,
	)
}

func UserResetPasswordToUserResetPasswordRepositoryMapper(userResetPassword userModel.UserResetPassword) UserResetPasswordRepository {
	return NewUserResetPasswordRepository(
		userResetPassword.Password,
	)
}

func UserRepositoryToUsersRepositoryMapper(usersRepository []UserRepository) UsersRepository {
	return NewUsersRepository(
		append([]UserRepository{}, usersRepository...),
	)
}
